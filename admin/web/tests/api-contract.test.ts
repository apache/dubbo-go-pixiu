import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { pixiuAdminApi } from '../src/api'
import { resourceApi } from '../src/services/resource-api'

beforeEach(() => {
  const storage = new Map<string, string>()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => storage.get(key) ?? null,
    setItem: (key: string, value: string) => storage.set(key, value),
    removeItem: (key: string) => storage.delete(key),
  })
})

afterEach(() => vi.unstubAllGlobals())

describe('Pixiu Admin API contract', () => {
  it('keeps resource and method routes aligned with the backend', () => {
    expect(pixiuAdminApi.resources).toEqual({
      list: '/config/api/resource/list',
      detail: '/config/api/resource/detail',
      create: '/config/api/resource',
      update: '/config/api/resource',
      remove: '/config/api/resource',
    })
    expect(pixiuAdminApi.methods).toEqual({
      list: '/config/api/resource/method/list',
      detail: '/config/api/resource/method/detail',
      create: '/config/api/resource/method',
      update: '/config/api/resource/method',
      remove: '/config/api/resource/method',
    })
  })

  it('does not reintroduce the deprecated api-admin proxy prefix', () => {
    expect(JSON.stringify(pixiuAdminApi)).not.toContain('/api-admin')
  })

  it('uses the shared contract for runtime resource requests', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify({ code: '10001', data: [] }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await resourceApi.list()

    expect(fetchMock).toHaveBeenCalledWith(
      pixiuAdminApi.resources.list,
      expect.objectContaining({ headers: expect.any(Headers) }),
    )
  })

  it('surfaces malformed list responses instead of treating them as empty', async () => {
    vi.stubGlobal(
      'fetch',
      vi
        .fn()
        .mockResolvedValue(
          new Response(JSON.stringify({ code: '10001', data: 'not-an-array' }), { status: 200 }),
        ),
    )

    await expect(resourceApi.list()).rejects.toMatchObject({ code: 'BAD_RESPONSE' })
  })

  it('does not accept an error HTTP status as a successful API response', async () => {
    vi.stubGlobal(
      'fetch',
      vi
        .fn()
        .mockResolvedValue(
          new Response(JSON.stringify({ message: 'backend unavailable' }), { status: 503 }),
        ),
    )

    await expect(resourceApi.list()).rejects.toMatchObject({ code: 'HTTP_503' })
  })
})
