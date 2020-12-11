// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import "cmd/compile/internal/types"

// select
func typecheckselect(sel *Node) {
	var def *Node
	lno := setlineno(sel)
	typecheckslice(sel.Ninit.Slice(), ctxStmt)
	for _, ncase := range sel.List.Slice() {
		if ncase.Op != OCASE {
			setlineno(ncase)
			Fatalf("typecheckselect %v", ncase.Op)
		}

		if ncase.List.Len() == 0 {
			// default
			if def != nil {
				yyerrorl(ncase.Pos, "multiple defaults in select (first at %v)", def.Line())
			} else {
				def = ncase
			}
		} else if ncase.List.Len() > 1 {
			yyerrorl(ncase.Pos, "select cases cannot be lists")
		} else {
			ncase.List.SetFirst(typecheck(ncase.List.First(), ctxStmt))
			n := ncase.List.First()
			ncase.Left = n
			ncase.List.Set(nil)
			switch n.Op {
			default:
				pos := n.Pos
				if n.Op == ONAME {
					// We don't have the right position for ONAME nodes (see #15459 and
					// others). Using ncase.Pos for now as it will provide the correct
					// line number (assuming the expression follows the "case" keyword
					// on the same line). This matches the approach before 1.10.
					pos = ncase.Pos
				}
				yyerrorl(pos, "select case must be receive, send or assign recv")

			// convert x = <-c into OSELRECV(x, <-c).
			// remove implicit conversions; the eventual assignment
			// will reintroduce them.
			case OAS:
				if (n.Right.Op == OCONVNOP || n.Right.Op == OCONVIFACE) && n.Right.Implicit() {
					n.Right = n.Right.Left
				}

				if n.Right.Op != ORECV {
					yyerrorl(n.Pos, "select assignment must have receive on right hand side")
					break
				}

				n.Op = OSELRECV

				// convert x, ok = <-c into OSELRECV2(x, <-c) with ntest=ok
			case OAS2RECV:
				if n.Right.Op != ORECV {
					yyerrorl(n.Pos, "select assignment must have receive on right hand side")
					break
				}

				n.Op = OSELRECV2
				n.Left = n.List.First()
				n.List.Set1(n.List.Second())

				// convert <-c into OSELRECV(N, <-c)
			case ORECV:
				n = nodl(n.Pos, OSELRECV, nil, n)

				n.SetTypecheck(1)
				ncase.Left = n

			case OSEND:
				break
			}
		}

		typecheckslice(ncase.Nbody.Slice(), ctxStmt)
	}

	lineno = lno
}

func walkselect(sel *Node) {
	lno := setlineno(sel)
	if sel.Nbody.Len() != 0 {
		Fatalf("double walkselect")
	}

	init := sel.Ninit.Slice()
	sel.Ninit.Set(nil)

	init = append(init, walkselectcases(&sel.List)...)
	sel.List.Set(nil)

	sel.Nbody.Set(init)
	walkstmtlist(sel.Nbody.Slice())

	lineno = lno
}

func walkselectcases(cases *Nodes) []*Node {
	n := cases.Len()
	sellineno := lineno

	// optimization: zero-case select
	if n == 0 {
		return []*Node{mkcall("block", nil, nil)}
	}

	// optimization: one-case select: single op.
	// TODO(rsc): Reenable optimization once order.go can handle it.
	// golang.org/issue/7672.
	if n == 1 {
		cas := cases.First()
		setlineno(cas)
		l := cas.Ninit.Slice()
		if cas.Left != nil { // not default:
			n := cas.Left
			l = append(l, n.Ninit.Slice()...)
			n.Ninit.Set(nil)
			var ch *Node
			switch n.Op {
			default:
				Fatalf("select %v", n.Op)

				// ok already
			case OSEND:
				ch = n.Left

			case OSELRECV, OSELRECV2:
				ch = n.Right.Left
				if n.Op == OSELRECV || n.List.Len() == 0 {
					if n.Left == nil {
						n = n.Right
					} else {
						n.Op = OAS
					}
					break
				}

				if n.Left == nil {
					nblank = typecheck(nblank, ctxExpr|ctxAssign)
					n.Left = nblank
				}

				n.Op = OAS2
				n.List.Prepend(n.Left)
				n.Rlist.Set1(n.Right)
				n.Right = nil
				n.Left = nil
				n.SetTypecheck(0)
				n = typecheck(n, ctxStmt)
			}

			// if ch == nil { block() }; n;
			a := nod(OIF, nil, nil)

			a.Left = nod(OEQ, ch, nodnil())
			var ln Nodes
			ln.Set(l)
			a.Nbody.Set1(mkcall("block", nil, &ln))
			l = ln.Slice()
			a = typecheck(a, ctxStmt)
			l = append(l, a, n)
		}

		l = append(l, cas.Nbody.Slice()...)
		l = append(l, nod(OBREAK, nil, nil))
		return l
	}

	// convert case value arguments to addresses.
	// this rewrite is used by both the general code and the next optimization.
	for _, cas := range cases.Slice() {
		setlineno(cas)
		n := cas.Left
		if n == nil {
			continue
		}
		switch n.Op {
		case OSEND:
			n.Right = nod(OADDR, n.Right, nil)
			n.Right = typecheck(n.Right, ctxExpr)

		case OSELRECV, OSELRECV2:
			if n.Op == OSELRECV2 && n.List.Len() == 0 {
				n.Op = OSELRECV
			}

			if n.Left != nil {
				n.Left = nod(OADDR, n.Left, nil)
				n.Left = typecheck(n.Left, ctxExpr)
			}
		}
	}

	// optimization: two-case select but one is default: single non-blocking op.
	if n == 2 && (cases.First().Left == nil || cases.Second().Left == nil) {
		var cas *Node
		var dflt *Node
		if cases.First().Left == nil {
			cas = cases.Second()
			dflt = cases.First()
		} else {
			dflt = cases.Second()
			cas = cases.First()
		}

		n := cas.Left
		setlineno(n)
		r := nod(OIF, nil, nil)
		r.Ninit.Set(cas.Ninit.Slice())
		switch n.Op {
		default:
			Fatalf("select %v", n.Op)

		case OSEND:
			// if selectnbsend(c, v) { body } else { default body }
			ch := n.Left
			r.Left = mkcall1(chanfn("selectnbsend", 2, ch.Type), types.Types[TBOOL], &r.Ninit, ch, n.Right)

		case OSELRECV:
			// if selectnbrecv(&v, c) { body } else { default body }
			r = nod(OIF, nil, nil)
			r.Ninit.Set(cas.Ninit.Slice())
			ch := n.Right.Left
			elem := n.Left
			if elem == nil {
				elem = nodnil()
			}
			r.Left = mkcall1(chanfn("selectnbrecv", 2, ch.Type), types.Types[TBOOL], &r.Ninit, elem, ch)

		case OSELRECV2:
			// if selectnbrecv2(&v, &received, c) { body } else { default body }
			r = nod(OIF, nil, nil)
			r.Ninit.Set(cas.Ninit.Slice())
			ch := n.Right.Left
			elem := n.Left
			if elem == nil {
				elem = nodnil()
			}
			receivedp := nod(OADDR, n.List.First(), nil)
			receivedp = typecheck(receivedp, ctxExpr)
			r.Left = mkcall1(chanfn("selectnbrecv2", 2, ch.Type), types.Types[TBOOL], &r.Ninit, elem, receivedp, ch)
		}

		r.Left = typecheck(r.Left, ctxExpr)
		r.Nbody.Set(cas.Nbody.Slice())
		r.Rlist.Set(append(dflt.Ninit.Slice(), dflt.Nbody.Slice()...))
		return []*Node{r, nod(OBREAK, nil, nil)}
	}

	var init []*Node

	// generate sel-struct
	lineno = sellineno
	selv := temp(types.NewArray(scasetype(), int64(n)))
	r := nod(OAS, selv, nil)
	r = typecheck(r, ctxStmt)
	init = append(init, r)

	order := temp(types.NewArray(types.Types[TUINT16], 2*int64(n)))
	r = nod(OAS, order, nil)
	r = typecheck(r, ctxStmt)
	init = append(init, r)

	// register cases
	for i, cas := range cases.Slice() {
		setlineno(cas)

		init = append(init, cas.Ninit.Slice()...)
		cas.Ninit.Set(nil)

		// Keep in sync with runtime/select.go.
		const (
			caseNil = iota
			caseRecv
			caseSend
			caseDefault
		)

		var c, elem *Node
		var kind int64 = caseDefault

		if n := cas.Left; n != nil {
			init = append(init, n.Ninit.Slice()...)

			switch n.Op {
			default:
				Fatalf("select %v", n.Op)
			case OSEND:
				kind = caseSend
				c = n.Left
				elem = n.Right
			case OSELRECV, OSELRECV2:
				kind = caseRecv
				c = n.Right.Left
				elem = n.Left
			}
		}

		setField := func(f string, val *Node) {
			r := nod(OAS, nodSym(ODOT, nod(OINDEX, selv, nodintconst(int64(i))), lookup(f)), val)
			r = typecheck(r, ctxStmt)
			init = append(init, r)
		}

		setField("kind", nodintconst(kind))
		if c != nil {
			c = convnop(c, types.Types[TUNSAFEPTR])
			setField("c", c)
		}
		if elem != nil {
			elem = convnop(elem, types.Types[TUNSAFEPTR])
			setField("elem", elem)
		}

		// TODO(mdempsky): There should be a cleaner way to
		// handle this.
		if instrumenting {
			r = mkcall("selectsetpc", nil, nil, bytePtrToIndex(selv, int64(i)))
			init = append(init, r)
		}
	}

	// run the select
	lineno = sellineno
	chosen := temp(types.Types[TINT])
	recvOK := temp(types.Types[TBOOL])
	r = nod(OAS2, nil, nil)
	r.List.Set2(chosen, recvOK)
	fn := syslook("selectgo")
	r.Rlist.Set1(mkcall1(fn, fn.Type.Results(), nil, bytePtrToIndex(selv, 0), bytePtrToIndex(order, 0), nodintconst(int64(n))))
	r = typecheck(r, ctxStmt)
	init = append(init, r)

	// selv and order are no longer alive after selectgo.
	init = append(init, nod(OVARKILL, selv, nil))
	init = append(init, nod(OVARKILL, order, nil))

	// dispatch cases
	for i, cas := range cases.Slice() {
		setlineno(cas)

		cond := nod(OEQ, chosen, nodintconst(int64(i)))
		cond = typecheck(cond, ctxExpr)
		cond = defaultlit(cond, nil)

		r = nod(OIF, cond, nil)

		if n := cas.Left; n != nil && n.Op == OSELRECV2 {
			x := nod(OAS, n.List.First(), recvOK)
			x = typecheck(x, ctxStmt)
			r.Nbody.Append(x)
		}

		r.Nbody.AppendNodes(&cas.Nbody)
		r.Nbody.Append(nod(OBREAK, nil, nil))
		init = append(init, r)
	}

	return init
}

// bytePtrToIndex returns a Node representing "(*byte)(&n[i])".
func bytePtrToIndex(n *Node, i int64) *Node {
	s := nod(OADDR, nod(OINDEX, n, nodintconst(i)), nil)
	t := types.NewPtr(types.Types[TUINT8])
	return convnop(s, t)
}

var scase *types.Type

// Keep in sync with src/runtime/select.go.
func scasetype() *types.Type {
	if scase == nil {
		scase = tostruct([]*Node{
			namedfield("c", types.Types[TUNSAFEPTR]),
			namedfield("elem", types.Types[TUNSAFEPTR]),
			namedfield("kind", types.Types[TUINT16]),
			namedfield("pc", types.Types[TUINTPTR]),
			namedfield("releasetime", types.Types[TINT64]),
		})
		scase.SetNoalg(true)
	}
	return scase
}
