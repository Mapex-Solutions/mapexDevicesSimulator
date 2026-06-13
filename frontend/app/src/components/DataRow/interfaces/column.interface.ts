import type { DataRowColumn } from './dataRow.interface';

/** Props shared by every DataRow column renderer. */
export interface DataRowColumnProps {
  value: unknown;
  column: DataRowColumn;
  row: Record<string, unknown>;
  mobile?: boolean;
}
