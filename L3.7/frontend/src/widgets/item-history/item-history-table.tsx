import { useItemHistory } from '@/entities/item/queries';
import { itemApi } from '@/entities/item/api';
import { Button } from '@/shared/ui/button';
import dayjs from 'dayjs';

function toRecord(value: unknown): Record<string, unknown> | null {
  if (!value || typeof value !== 'object') return null;
  if (Array.isArray(value)) return null;
  return value as Record<string, unknown>;
}

function diffFields(oldObj: unknown, newObj: unknown): Array<{ field: string; from?: unknown; to?: unknown }> {
  const oldR = toRecord(oldObj) || {};
  const newR = toRecord(newObj) || {};

  const keys = new Set<string>([...Object.keys(oldR), ...Object.keys(newR)]);
  const out: Array<{ field: string; from?: unknown; to?: unknown }> = [];

  for (const k of Array.from(keys).sort()) {
    const a = oldR[k];
    const b = newR[k];
    if (JSON.stringify(a) !== JSON.stringify(b)) {
      out.push({ field: k, from: a, to: b });
    }
  }

  return out;
}

function renderValue(v: unknown): string {
  if (v == null) return '';
  if (typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean') return String(v);
  try {
    return JSON.stringify(v);
  } catch {
    return String(v);
  }
}

export const ItemHistoryTable = ({ itemId }: { itemId: string }) => {
  const { data, isLoading, isError, error } = useItemHistory(itemId);

  if (isLoading) {
    return <div className="text-sm text-muted-foreground">Loading history…</div>;
  }

  if (isError) {
    return (
      <div className="text-sm text-destructive">
        Failed to load history: {(error as Error)?.message || 'unknown error'}
      </div>
    );
  }

  const rows = Array.isArray(data) ? data : [];

  const downloadCSV = async () => {
    const { blob, filename } = await itemApi.downloadHistoryCSV(itemId);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.style.display = 'none';
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div className="text-sm text-muted-foreground">{rows.length} record(s)</div>
        <Button variant="outline" size="sm" onClick={downloadCSV}>
          Export CSV
        </Button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="text-left text-muted-foreground">
            <tr className="border-b">
              <th className="py-2 pr-4">When</th>
              <th className="py-2 pr-4">Action</th>
              <th className="py-2 pr-4">User</th>
              <th className="py-2 pr-4">Role</th>
              <th className="py-2 pr-4">Changes</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.id} className="border-b align-top">
                <td className="py-2 pr-4 whitespace-nowrap">
                  {dayjs(r.created_at).format('YYYY-MM-DD HH:mm')}
                </td>
                <td className="py-2 pr-4 font-mono">{r.action}</td>
                <td className="py-2 pr-4">{r.username || '-'}</td>
                <td className="py-2 pr-4">{r.role || '-'}</td>
                <td className="py-2 pr-4">
                  <div className="space-y-1">
                    {diffFields(r.old_data, r.new_data).map((d) => (
                      <div key={d.field} className="text-xs">
                        <span className="font-mono">{d.field}</span>
                        {r.action === 'INSERT' ? (
                          <span className="text-muted-foreground"> = {renderValue(d.to)}</span>
                        ) : r.action === 'DELETE' ? (
                          <span className="text-muted-foreground"> was {renderValue(d.from)}</span>
                        ) : (
                          <span className="text-muted-foreground">
                            : {renderValue(d.from)} → {renderValue(d.to)}
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};




