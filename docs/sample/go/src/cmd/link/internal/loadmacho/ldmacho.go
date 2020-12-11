// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package loadmacho implements a Mach-O file reader.
package loadmacho

import (
	"bytes"
	"cmd/internal/bio"
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"cmd/link/internal/loader"
	"cmd/link/internal/sym"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

/*
Derived from Plan 9 from User Space's src/libmach/elf.h, elf.c
http://code.swtch.com/plan9port/src/tip/src/libmach/

	Copyright © 2004 Russ Cox.
	Portions Copyright © 2008-2010 Google Inc.
	Portions Copyright © 2010 The Go Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// TODO(crawshaw): de-duplicate these symbols with cmd/internal/ld
const (
	MACHO_X86_64_RELOC_UNSIGNED = 0
	MACHO_X86_64_RELOC_SIGNED   = 1
	MACHO_FAKE_GOTPCREL         = 100
)

type ldMachoObj struct {
	f          *bio.Reader
	base       int64 // off in f where Mach-O begins
	length     int64 // length of Mach-O
	is64       bool
	name       string
	e          binary.ByteOrder
	cputype    uint
	subcputype uint
	filetype   uint32
	flags      uint32
	cmd        []ldMachoCmd
	ncmd       uint
}

type ldMachoCmd struct {
	type_ int
	off   uint32
	size  uint32
	seg   ldMachoSeg
	sym   ldMachoSymtab
	dsym  ldMachoDysymtab
}

type ldMachoSeg struct {
	name     string
	vmaddr   uint64
	vmsize   uint64
	fileoff  uint32
	filesz   uint32
	maxprot  uint32
	initprot uint32
	nsect    uint32
	flags    uint32
	sect     []ldMachoSect
}

type ldMachoSect struct {
	name    string
	segname string
	addr    uint64
	size    uint64
	off     uint32
	align   uint32
	reloff  uint32
	nreloc  uint32
	flags   uint32
	res1    uint32
	res2    uint32
	sym     *sym.Symbol
	rel     []ldMachoRel
}

type ldMachoRel struct {
	addr      uint32
	symnum    uint32
	pcrel     uint8
	length    uint8
	extrn     uint8
	type_     uint8
	scattered uint8
	value     uint32
}

type ldMachoSymtab struct {
	symoff  uint32
	nsym    uint32
	stroff  uint32
	strsize uint32
	str     []byte
	sym     []ldMachoSym
}

type ldMachoSym struct {
	name    string
	type_   uint8
	sectnum uint8
	desc    uint16
	kind    int8
	value   uint64
	sym     *sym.Symbol
}

type ldMachoDysymtab struct {
	ilocalsym      uint32
	nlocalsym      uint32
	iextdefsym     uint32
	nextdefsym     uint32
	iundefsym      uint32
	nundefsym      uint32
	tocoff         uint32
	ntoc           uint32
	modtaboff      uint32
	nmodtab        uint32
	extrefsymoff   uint32
	nextrefsyms    uint32
	indirectsymoff uint32
	nindirectsyms  uint32
	extreloff      uint32
	nextrel        uint32
	locreloff      uint32
	nlocrel        uint32
	indir          []uint32
}

// ldMachoSym.type_
const (
	N_EXT  = 0x01
	N_TYPE = 0x1e
	N_STAB = 0xe0
)

// ldMachoSym.desc
const (
	N_WEAK_REF = 0x40
	N_WEAK_DEF = 0x80
)

const (
	LdMachoCpuVax         = 1
	LdMachoCpu68000       = 6
	LdMachoCpu386         = 7
	LdMachoCpuAmd64       = 0x1000007
	LdMachoCpuMips        = 8
	LdMachoCpu98000       = 10
	LdMachoCpuHppa        = 11
	LdMachoCpuArm         = 12
	LdMachoCpu88000       = 13
	LdMachoCpuSparc       = 14
	LdMachoCpu860         = 15
	LdMachoCpuAlpha       = 16
	LdMachoCpuPower       = 18
	LdMachoCmdSegment     = 1
	LdMachoCmdSymtab      = 2
	LdMachoCmdSymseg      = 3
	LdMachoCmdThread      = 4
	LdMachoCmdDysymtab    = 11
	LdMachoCmdSegment64   = 25
	LdMachoFileObject     = 1
	LdMachoFileExecutable = 2
	LdMachoFileFvmlib     = 3
	LdMachoFileCore       = 4
	LdMachoFilePreload    = 5
)

func unpackcmd(p []byte, m *ldMachoObj, c *ldMachoCmd, type_ uint, sz uint) int {
	e4 := m.e.Uint32
	e8 := m.e.Uint64

	c.type_ = int(type_)
	c.size = uint32(sz)
	switch type_ {
	default:
		return -1

	case LdMachoCmdSegment:
		if sz < 56 {
			return -1
		}
		c.seg.name = cstring(p[8:24])
		c.seg.vmaddr = uint64(e4(p[24:]))
		c.seg.vmsize = uint64(e4(p[28:]))
		c.seg.fileoff = e4(p[32:])
		c.seg.filesz = e4(p[36:])
		c.seg.maxprot = e4(p[40:])
		c.seg.initprot = e4(p[44:])
		c.seg.nsect = e4(p[48:])
		c.seg.flags = e4(p[52:])
		c.seg.sect = make([]ldMachoSect, c.seg.nsect)
		if uint32(sz) < 56+c.seg.nsect*68 {
			return -1
		}
		p = p[56:]
		var s *ldMachoSect
		for i := 0; uint32(i) < c.seg.nsect; i++ {
			s = &c.seg.sect[i]
			s.name = cstring(p[0:16])
			s.segname = cstring(p[16:32])
			s.addr = uint64(e4(p[32:]))
			s.size = uint64(e4(p[36:]))
			s.off = e4(p[40:])
			s.align = e4(p[44:])
			s.reloff = e4(p[48:])
			s.nreloc = e4(p[52:])
			s.flags = e4(p[56:])
			s.res1 = e4(p[60:])
			s.res2 = e4(p[64:])
			p = p[68:]
		}

	case LdMachoCmdSegment64:
		if sz < 72 {
			return -1
		}
		c.seg.name = cstring(p[8:24])
		c.seg.vmaddr = e8(p[24:])
		c.seg.vmsize = e8(p[32:])
		c.seg.fileoff = uint32(e8(p[40:]))
		c.seg.filesz = uint32(e8(p[48:]))
		c.seg.maxprot = e4(p[56:])
		c.seg.initprot = e4(p[60:])
		c.seg.nsect = e4(p[64:])
		c.seg.flags = e4(p[68:])
		c.seg.sect = make([]ldMachoSect, c.seg.nsect)
		if uint32(sz) < 72+c.seg.nsect*80 {
			return -1
		}
		p = p[72:]
		var s *ldMachoSect
		for i := 0; uint32(i) < c.seg.nsect; i++ {
			s = &c.seg.sect[i]
			s.name = cstring(p[0:16])
			s.segname = cstring(p[16:32])
			s.addr = e8(p[32:])
			s.size = e8(p[40:])
			s.off = e4(p[48:])
			s.align = e4(p[52:])
			s.reloff = e4(p[56:])
			s.nreloc = e4(p[60:])
			s.flags = e4(p[64:])
			s.res1 = e4(p[68:])
			s.res2 = e4(p[72:])

			// p+76 is reserved
			p = p[80:]
		}

	case LdMachoCmdSymtab:
		if sz < 24 {
			return -1
		}
		c.sym.symoff = e4(p[8:])
		c.sym.nsym = e4(p[12:])
		c.sym.stroff = e4(p[16:])
		c.sym.strsize = e4(p[20:])

	case LdMachoCmdDysymtab:
		if sz < 80 {
			return -1
		}
		c.dsym.ilocalsym = e4(p[8:])
		c.dsym.nlocalsym = e4(p[12:])
		c.dsym.iextdefsym = e4(p[16:])
		c.dsym.nextdefsym = e4(p[20:])
		c.dsym.iundefsym = e4(p[24:])
		c.dsym.nundefsym = e4(p[28:])
		c.dsym.tocoff = e4(p[32:])
		c.dsym.ntoc = e4(p[36:])
		c.dsym.modtaboff = e4(p[40:])
		c.dsym.nmodtab = e4(p[44:])
		c.dsym.extrefsymoff = e4(p[48:])
		c.dsym.nextrefsyms = e4(p[52:])
		c.dsym.indirectsymoff = e4(p[56:])
		c.dsym.nindirectsyms = e4(p[60:])
		c.dsym.extreloff = e4(p[64:])
		c.dsym.nextrel = e4(p[68:])
		c.dsym.locreloff = e4(p[72:])
		c.dsym.nlocrel = e4(p[76:])
	}

	return 0
}

func macholoadrel(m *ldMachoObj, sect *ldMachoSect) int {
	if sect.rel != nil || sect.nreloc == 0 {
		return 0
	}
	rel := make([]ldMachoRel, sect.nreloc)
	n := int(sect.nreloc * 8)
	buf := make([]byte, n)
	m.f.MustSeek(m.base+int64(sect.reloff), 0)
	if _, err := io.ReadFull(m.f, buf); err != nil {
		return -1
	}
	for i := uint32(0); i < sect.nreloc; i++ {
		r := &rel[i]
		p := buf[i*8:]
		r.addr = m.e.Uint32(p)

		// TODO(rsc): Wrong interpretation for big-endian bitfields?
		if r.addr&0x80000000 != 0 {
			// scatterbrained relocation
			r.scattered = 1

			v := r.addr >> 24
			r.addr &= 0xFFFFFF
			r.type_ = uint8(v & 0xF)
			v >>= 4
			r.length = 1 << (v & 3)
			v >>= 2
			r.pcrel = uint8(v & 1)
			r.value = m.e.Uint32(p[4:])
		} else {
			v := m.e.Uint32(p[4:])
			r.symnum = v & 0xFFFFFF
			v >>= 24
			r.pcrel = uint8(v & 1)
			v >>= 1
			r.length = 1 << (v & 3)
			v >>= 2
			r.extrn = uint8(v & 1)
			v >>= 1
			r.type_ = uint8(v)
		}
	}

	sect.rel = rel
	return 0
}

func macholoaddsym(m *ldMachoObj, d *ldMachoDysymtab) int {
	n := int(d.nindirectsyms)

	p := make([]byte, n*4)
	m.f.MustSeek(m.base+int64(d.indirectsymoff), 0)
	if _, err := io.ReadFull(m.f, p); err != nil {
		return -1
	}

	d.indir = make([]uint32, n)
	for i := 0; i < n; i++ {
		d.indir[i] = m.e.Uint32(p[4*i:])
	}
	return 0
}

func macholoadsym(m *ldMachoObj, symtab *ldMachoSymtab) int {
	if symtab.sym != nil {
		return 0
	}

	strbuf := make([]byte, symtab.strsize)
	m.f.MustSeek(m.base+int64(symtab.stroff), 0)
	if _, err := io.ReadFull(m.f, strbuf); err != nil {
		return -1
	}

	symsize := 12
	if m.is64 {
		symsize = 16
	}
	n := int(symtab.nsym * uint32(symsize))
	symbuf := make([]byte, n)
	m.f.MustSeek(m.base+int64(symtab.symoff), 0)
	if _, err := io.ReadFull(m.f, symbuf); err != nil {
		return -1
	}
	sym := make([]ldMachoSym, symtab.nsym)
	p := symbuf
	for i := uint32(0); i < symtab.nsym; i++ {
		s := &sym[i]
		v := m.e.Uint32(p)
		if v >= symtab.strsize {
			return -1
		}
		s.name = cstring(strbuf[v:])
		s.type_ = p[4]
		s.sectnum = p[5]
		s.desc = m.e.Uint16(p[6:])
		if m.is64 {
			s.value = m.e.Uint64(p[8:])
		} else {
			s.value = uint64(m.e.Uint32(p[8:]))
		}
		p = p[symsize:]
	}

	symtab.str = strbuf
	symtab.sym = sym
	return 0
}

func Load(l *loader.Loader, arch *sys.Arch, syms *sym.Symbols, f *bio.Reader, pkg string, length int64, pn string) ([]*sym.Symbol, error) {
	newSym := func(name string, version int) *sym.Symbol {
		return l.LookupOrCreate(name, version, syms)
	}
	return load(arch, syms.IncVersion(), newSym, f, pkg, length, pn)
}

func LoadOld(arch *sys.Arch, syms *sym.Symbols, f *bio.Reader, pkg string, length int64, pn string) (textp []*sym.Symbol, err error) {
	return load(arch, syms.IncVersion(), syms.Lookup, f, pkg, length, pn)
}

// load the Mach-O file pn from f.
// Symbols are written into syms, and a slice of the text symbols is returned.
func load(arch *sys.Arch, localSymVersion int, lookup func(string, int) *sym.Symbol, f *bio.Reader, pkg string, length int64, pn string) (textp []*sym.Symbol, err error) {
	errorf := func(str string, args ...interface{}) ([]*sym.Symbol, error) {
		return nil, fmt.Errorf("loadmacho: %v: %v", pn, fmt.Sprintf(str, args...))
	}

	base := f.Offset()

	var hdr [7 * 4]uint8
	if _, err := io.ReadFull(f, hdr[:]); err != nil {
		return errorf("reading hdr: %v", err)
	}

	var e binary.ByteOrder
	if binary.BigEndian.Uint32(hdr[:])&^1 == 0xFEEDFACE {
		e = binary.BigEndian
	} else if binary.LittleEndian.Uint32(hdr[:])&^1 == 0xFEEDFACE {
		e = binary.LittleEndian
	} else {
		return errorf("bad magic - not mach-o file")
	}

	is64 := e.Uint32(hdr[:]) == 0xFEEDFACF
	ncmd := e.Uint32(hdr[4*4:])
	cmdsz := e.Uint32(hdr[5*4:])
	if ncmd > 0x10000 || cmdsz >= 0x01000000 {
		return errorf("implausible mach-o header ncmd=%d cmdsz=%d", ncmd, cmdsz)
	}

	if is64 {
		f.MustSeek(4, 1) // skip reserved word in header
	}

	m := &ldMachoObj{
		f:          f,
		e:          e,
		cputype:    uint(e.Uint32(hdr[1*4:])),
		subcputype: uint(e.Uint32(hdr[2*4:])),
		filetype:   e.Uint32(hdr[3*4:]),
		ncmd:       uint(ncmd),
		flags:      e.Uint32(hdr[6*4:]),
		is64:       is64,
		base:       base,
		length:     length,
		name:       pn,
	}

	switch arch.Family {
	default:
		return errorf("mach-o %s unimplemented", arch.Name)

	case sys.AMD64:
		if e != binary.LittleEndian || m.cputype != LdMachoCpuAmd64 {
			return errorf("mach-o object but not amd64")
		}

	case sys.I386:
		if e != binary.LittleEndian || m.cputype != LdMachoCpu386 {
			return errorf("mach-o object but not 386")
		}
	}

	m.cmd = make([]ldMachoCmd, ncmd)
	cmdp := make([]byte, cmdsz)
	if _, err := io.ReadFull(f, cmdp); err != nil {
		return errorf("reading cmds: %v", err)
	}

	// read and parse load commands
	var c *ldMachoCmd

	var symtab *ldMachoSymtab
	var dsymtab *ldMachoDysymtab

	off := uint32(len(hdr))
	for i := uint32(0); i < ncmd; i++ {
		ty := e.Uint32(cmdp)
		sz := e.Uint32(cmdp[4:])
		m.cmd[i].off = off
		unpackcmd(cmdp, m, &m.cmd[i], uint(ty), uint(sz))
		cmdp = cmdp[sz:]
		off += sz
		if ty == LdMachoCmdSymtab {
			if symtab != nil {
				return errorf("multiple symbol tables")
			}

			symtab = &m.cmd[i].sym
			macholoadsym(m, symtab)
		}

		if ty == LdMachoCmdDysymtab {
			dsymtab = &m.cmd[i].dsym
			macholoaddsym(m, dsymtab)
		}

		if (is64 && ty == LdMachoCmdSegment64) || (!is64 && ty == LdMachoCmdSegment) {
			if c != nil {
				return errorf("multiple load commands")
			}

			c = &m.cmd[i]
		}
	}

	// load text and data segments into memory.
	// they are not as small as the load commands, but we'll need
	// the memory anyway for the symbol images, so we might
	// as well use one large chunk.
	if c == nil {
		return errorf("no load command")
	}

	if symtab == nil {
		// our work is done here - no symbols means nothing can refer to this file
		return
	}

	if int64(c.seg.fileoff+c.seg.filesz) >= length {
		return errorf("load segment out of range")
	}

	f.MustSeek(m.base+int64(c.seg.fileoff), 0)
	dat := make([]byte, c.seg.filesz)
	if _, err := io.ReadFull(f, dat); err != nil {
		return errorf("cannot load object data: %v", err)
	}

	for i := uint32(0); i < c.seg.nsect; i++ {
		sect := &c.seg.sect[i]
		if sect.segname != "__TEXT" && sect.segname != "__DATA" {
			continue
		}
		if sect.name == "__eh_frame" {
			continue
		}
		name := fmt.Sprintf("%s(%s/%s)", pkg, sect.segname, sect.name)
		s := lookup(name, localSymVersion)
		if s.Type != 0 {
			return errorf("duplicate %s/%s", sect.segname, sect.name)
		}

		if sect.flags&0xff == 1 { // S_ZEROFILL
			s.P = make([]byte, sect.size)
		} else {
			s.P = dat[sect.addr-c.seg.vmaddr:][:sect.size]
		}
		s.Size = int64(len(s.P))

		if sect.segname == "__TEXT" {
			if sect.name == "__text" {
				s.Type = sym.STEXT
			} else {
				s.Type = sym.SRODATA
			}
		} else {
			if sect.name == "__bss" {
				s.Type = sym.SNOPTRBSS
				s.P = s.P[:0]
			} else {
				s.Type = sym.SNOPTRDATA
			}
		}

		sect.sym = s
	}

	// enter sub-symbols into symbol table.
	// have to guess sizes from next symbol.
	for i := uint32(0); i < symtab.nsym; i++ {
		machsym := &symtab.sym[i]
		if machsym.type_&N_STAB != 0 {
			continue
		}

		// TODO: check sym->type against outer->type.
		name := machsym.name

		if name[0] == '_' && name[1] != '\x00' {
			name = name[1:]
		}
		v := 0
		if machsym.type_&N_EXT == 0 {
			v = localSymVersion
		}
		s := lookup(name, v)
		if machsym.type_&N_EXT == 0 {
			s.Attr |= sym.AttrDuplicateOK
		}
		if machsym.desc&(N_WEAK_REF|N_WEAK_DEF) != 0 {
			s.Attr |= sym.AttrDuplicateOK
		}
		machsym.sym = s
		if machsym.sectnum == 0 { // undefined
			continue
		}
		if uint32(machsym.sectnum) > c.seg.nsect {
			return errorf("reference to invalid section %d", machsym.sectnum)
		}

		sect := &c.seg.sect[machsym.sectnum-1]
		outer := sect.sym
		if outer == nil {
			continue // ignore reference to invalid section
		}

		if s.Outer != nil {
			if s.Attr.DuplicateOK() {
				continue
			}
			return errorf("duplicate symbol reference: %s in both %s and %s", s.Name, s.Outer.Name, sect.sym.Name)
		}

		s.Type = outer.Type
		s.Attr |= sym.AttrSubSymbol
		s.Sub = outer.Sub
		outer.Sub = s
		s.Outer = outer
		s.Value = int64(machsym.value - sect.addr)
		if !s.Attr.CgoExportDynamic() {
			s.SetDynimplib("") // satisfy dynimport
		}
		if outer.Type == sym.STEXT {
			if s.Attr.External() && !s.Attr.DuplicateOK() {
				return errorf("%v: duplicate symbol definition", s)
			}
			s.Attr |= sym.AttrExternal
		}

		machsym.sym = s
	}

	// Sort outer lists by address, adding to textp.
	// This keeps textp in increasing address order.
	for i := 0; uint32(i) < c.seg.nsect; i++ {
		sect := &c.seg.sect[i]
		s := sect.sym
		if s == nil {
			continue
		}
		if s.Sub != nil {
			s.Sub = sym.SortSub(s.Sub)

			// assign sizes, now that we know symbols in sorted order.
			for s1 := s.Sub; s1 != nil; s1 = s1.Sub {
				if s1.Sub != nil {
					s1.Size = s1.Sub.Value - s1.Value
				} else {
					s1.Size = s.Value + s.Size - s1.Value
				}
			}
		}

		if s.Type == sym.STEXT {
			if s.Attr.OnList() {
				return errorf("symbol %s listed multiple times", s.Name)
			}
			s.Attr |= sym.AttrOnList
			textp = append(textp, s)
			for s1 := s.Sub; s1 != nil; s1 = s1.Sub {
				if s1.Attr.OnList() {
					return errorf("symbol %s listed multiple times", s1.Name)
				}
				s1.Attr |= sym.AttrOnList
				textp = append(textp, s1)
			}
		}
	}

	// load relocations
	for i := 0; uint32(i) < c.seg.nsect; i++ {
		sect := &c.seg.sect[i]
		s := sect.sym
		if s == nil {
			continue
		}
		macholoadrel(m, sect)
		if sect.rel == nil {
			continue
		}
		r := make([]sym.Reloc, sect.nreloc)
		rpi := 0
	Reloc:
		for j := uint32(0); j < sect.nreloc; j++ {
			rp := &r[rpi]
			rel := &sect.rel[j]
			if rel.scattered != 0 {
				if arch.Family != sys.I386 {
					// mach-o only uses scattered relocation on 32-bit platforms
					return errorf("%v: unexpected scattered relocation", s)
				}

				// on 386, rewrite scattered 4/1 relocation and some
				// scattered 2/1 relocation into the pseudo-pc-relative
				// reference that it is.
				// assume that the second in the pair is in this section
				// and use that as the pc-relative base.
				if j+1 >= sect.nreloc {
					return errorf("unsupported scattered relocation %d", int(rel.type_))
				}

				if sect.rel[j+1].scattered == 0 || sect.rel[j+1].type_ != 1 || (rel.type_ != 4 && rel.type_ != 2) || uint64(sect.rel[j+1].value) < sect.addr || uint64(sect.rel[j+1].value) >= sect.addr+sect.size {
					return errorf("unsupported scattered relocation %d/%d", int(rel.type_), int(sect.rel[j+1].type_))
				}

				rp.Siz = rel.length
				rp.Off = int32(rel.addr)

				// NOTE(rsc): I haven't worked out why (really when)
				// we should ignore the addend on a
				// scattered relocation, but it seems that the
				// common case is we ignore it.
				// It's likely that this is not strictly correct
				// and that the math should look something
				// like the non-scattered case below.
				rp.Add = 0

				// want to make it pc-relative aka relative to rp->off+4
				// but the scatter asks for relative to off = sect->rel[j+1].value - sect->addr.
				// adjust rp->add accordingly.
				rp.Type = objabi.R_PCREL

				rp.Add += int64(uint64(int64(rp.Off)+4) - (uint64(sect.rel[j+1].value) - sect.addr))

				// now consider the desired symbol.
				// find the section where it lives.
				for k := 0; uint32(k) < c.seg.nsect; k++ {
					ks := &c.seg.sect[k]
					if ks.addr <= uint64(rel.value) && uint64(rel.value) < ks.addr+ks.size {
						if ks.sym != nil {
							rp.Sym = ks.sym
							rp.Add += int64(uint64(rel.value) - ks.addr)
						} else if ks.segname == "__IMPORT" && ks.name == "__pointers" {
							// handle reference to __IMPORT/__pointers.
							// how much worse can this get?
							// why are we supporting 386 on the mac anyway?
							rp.Type = objabi.MachoRelocOffset + MACHO_FAKE_GOTPCREL

							// figure out which pointer this is a reference to.
							k = int(uint64(ks.res1) + (uint64(rel.value)-ks.addr)/4)

							// load indirect table for __pointers
							// fetch symbol number
							if dsymtab == nil || k < 0 || uint32(k) >= dsymtab.nindirectsyms || dsymtab.indir == nil {
								return errorf("invalid scattered relocation: indirect symbol reference out of range")
							}

							k = int(dsymtab.indir[k])
							if k < 0 || uint32(k) >= symtab.nsym {
								return errorf("invalid scattered relocation: symbol reference out of range")
							}

							rp.Sym = symtab.sym[k].sym
						} else {
							return errorf("unsupported scattered relocation: reference to %s/%s", ks.segname, ks.name)
						}

						rpi++

						// skip #1 of 2 rel; continue skips #2 of 2.
						j++

						continue Reloc
					}
				}

				return errorf("unsupported scattered relocation: invalid address %#x", rel.addr)
			}

			rp.Siz = rel.length
			rp.Type = objabi.MachoRelocOffset + (objabi.RelocType(rel.type_) << 1) + objabi.RelocType(rel.pcrel)
			rp.Off = int32(rel.addr)

			// Handle X86_64_RELOC_SIGNED referencing a section (rel->extrn == 0).
			if arch.Family == sys.AMD64 && rel.extrn == 0 && rel.type_ == MACHO_X86_64_RELOC_SIGNED {
				// Calculate the addend as the offset into the section.
				//
				// The rip-relative offset stored in the object file is encoded
				// as follows:
				//
				//    movsd	0x00000360(%rip),%xmm0
				//
				// To get the absolute address of the value this rip-relative address is pointing
				// to, we must add the address of the next instruction to it. This is done by
				// taking the address of the relocation and adding 4 to it (since the rip-relative
				// offset can at most be 32 bits long).  To calculate the offset into the section the
				// relocation is referencing, we subtract the vaddr of the start of the referenced
				// section found in the original object file.
				//
				// [For future reference, see Darwin's /usr/include/mach-o/x86_64/reloc.h]
				secaddr := c.seg.sect[rel.symnum-1].addr

				rp.Add = int64(uint64(int64(int32(e.Uint32(s.P[rp.Off:])))+int64(rp.Off)+4) - secaddr)
			} else {
				rp.Add = int64(int32(e.Uint32(s.P[rp.Off:])))
			}

			// An unsigned internal relocation has a value offset
			// by the section address.
			if arch.Family == sys.AMD64 && rel.extrn == 0 && rel.type_ == MACHO_X86_64_RELOC_UNSIGNED {
				secaddr := c.seg.sect[rel.symnum-1].addr
				rp.Add -= int64(secaddr)
			}

			// For i386 Mach-O PC-relative, the addend is written such that
			// it *is* the PC being subtracted. Use that to make
			// it match our version of PC-relative.
			if rel.pcrel != 0 && arch.Family == sys.I386 {
				rp.Add += int64(rp.Off) + int64(rp.Siz)
			}
			if rel.extrn == 0 {
				if rel.symnum < 1 || rel.symnum > c.seg.nsect {
					return errorf("invalid relocation: section reference out of range %d vs %d", rel.symnum, c.seg.nsect)
				}

				rp.Sym = c.seg.sect[rel.symnum-1].sym
				if rp.Sym == nil {
					return errorf("invalid relocation: %s", c.seg.sect[rel.symnum-1].name)
				}

				// References to symbols in other sections
				// include that information in the addend.
				// We only care about the delta from the
				// section base.
				if arch.Family == sys.I386 {
					rp.Add -= int64(c.seg.sect[rel.symnum-1].addr)
				}
			} else {
				if rel.symnum >= symtab.nsym {
					return errorf("invalid relocation: symbol reference out of range")
				}

				rp.Sym = symtab.sym[rel.symnum].sym
			}

			rpi++
		}

		sort.Sort(sym.RelocByOff(r[:rpi]))
		s.R = r
		s.R = s.R[:rpi]
	}

	return textp, nil
}

func cstring(x []byte) string {
	i := bytes.IndexByte(x, '\x00')
	if i >= 0 {
		x = x[:i]
	}
	return string(x)
}
