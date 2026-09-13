import React from 'react'
export function DataTable<T extends Record<string, unknown>>({
  data,
  columns,
}: {
  data: T[]
  columns: { key: keyof T; header: string }[]
}) {
  return (
    <table>
      <thead>
        <tr>
          {columns.map((c) => (
            <th key={String(c.key)}>{c.header}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((row, i) => (
          <tr key={i}>
            {columns.map((c) => (
              <td key={String(c.key)}>{String(row[c.key] ?? '')}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  )
}
