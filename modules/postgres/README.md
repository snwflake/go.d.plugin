<!--
title: "PostgreSQL monitoring with Netdata"
description: "Monitor connections, replication, databases, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/postgres/README.md
sidebar_label: "PostgresSQL"
-->

# PostgreSQL monitoring with Netdata

[PostgreSQL](https://www.postgresql.org/), also known as Postgres, is a free and open-source relational database
management system emphasizing extensibility and SQL compliance.

This module monitors one or more Postgres servers, depending on your configuration.

## Requirements

- PostgreSQL v9.4+
- User with granted `pg_monitor`
  or `pg_read_all_stat` [built-in role](https://www.postgresql.org/docs/current/predefined-roles.html).

## Metrics

All metrics have "postgres." prefix.

| Metric                               |        Scope        |                                                                 Dimensions                                                                 |     Units      |
|--------------------------------------|:-------------------:|:------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| connections_utilization              |       global        |                                                                    used                                                                    |   percentage   |
| connections_usage                    |       global        |                                                              available, used                                                               |  connections   |
| checkpoints                          |       global        |                                                            scheduled, requested                                                            | checkpoints/s  |
| checkpoint_time                      |       global        |                                                                write, sync                                                                 |  milliseconds  |
| bgwriter_buffers_alloc               |       global        |                                                                 allocated                                                                  |      B/s       |
| bgwriter_buffers_written             |       global        |                                                         checkpoint, backend, clean                                                         |      B/s       |
| bgwriter_maxwritten_clean            |       global        |                                                                 maxwritten                                                                 |    events/s    |
| bgwriter_buffers_backend_fsync       |       global        |                                                                   fsync                                                                    |  operations/s  |
| wal_writes                           |       global        |                                                                   writes                                                                   |      B/s       |
| wal_files                            |       global        |                                                             written, recycled                                                              |     files      |
| wal_archive_files                    |       global        |                                                                ready, done                                                                 |    files/s     |
| autovacuum_workers                   |       global        |                                       analyze, vacuum_analyze, vacuum, vacuum_freeze, brin_summarize                                       |    workers     |
| percent_towards_emergency_autovacuum |       global        |                                                            emergency_autovacuum                                                            |   percentage   |
| percent_towards_txid_wraparound      |       global        |                                                              txid_wraparound                                                               |   percentage   |
| oldest_transaction_xid               |       global        |                                                                    xid                                                                     |      xid       |
| catalog_relation_count               |       global        | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |   relations    |
| catalog_relation_size                |       global        | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |       B        |
| uptime                               |       global        |                                                                   uptime                                                                   |    seconds     |
| replication_standby_app_wal_delta    | standby application |                                             sent_delta, write_delta, flush_delta, replay_delta                                             |       B        |
| replication_standby_app_wal_lag      | standby application |                                                      write_lag, flush_lag, replay_lag                                                      |    seconds     |
| replication_standby_app_wal_lag      |  replication slot   |                                                        wal_keep, pg_replslot_files                                                         |     files      |
| db_transactions_ratio                |      database       |                                                            committed, rollback                                                             |   percentage   |
| db_transactions                      |      database       |                                                            committed, rollback                                                             | transactions/s |
| db_connections_utilization           |      database       |                                                                    used                                                                    |   percentage   |
| db_connections                       |      database       |                                                                connections                                                                 |  connections   |
| db_buffer_cache_hit_ratio            |      database       |                                                                 hit, miss                                                                  |   percentage   |
| db_blocks_read                       |      database       |                                                                memory, disk                                                                |    blocks/s    |
| db_rows_read_ratio                   |      database       |                                                             returned, fetched                                                              |   percentage   |
| db_rows_read                         |      database       |                                                             returned, fetched                                                              |     rows/s     |
| db_rows_written                      |      database       |                                                         inserted, deleted, updated                                                         |     rows/s     |
| db_conflicts                         |      database       |                                                                 conflicts                                                                  |   queries/s    |
| db_conflicts_stat                    |      database       |                                              tablespace, lock, snapshot, bufferpin, deadlock                                               |   queries/s    |
| db_deadlocks                         |      database       |                                                                 deadlocks                                                                  |  deadlocks/s   |
| db_locks_held                        |      database       |               access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive                |     locks      |
| db_locks_awaited                     |      database       |               access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive                |     locks      |
| db_temp_files                        |      database       |                                                                  written                                                                   |    files/s     |
| db_temp_files_data                   |      database       |                                                                  written                                                                   |      B/s       |
| db_size                              |      database       |                                                                    size                                                                    |       B        |

## Configuration

Edit the `go.d/postgres.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/postgres.conf
```

DSN (Data Source Name) may either be in URL format or key=word format.
See [Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for details.

```yaml
jobs:
  - name: local
    dsn: 'postgres://postgres:postgres@127.0.0.1:5432/postgres'

  - name: local
    dsn: 'host=/var/run/postgresql dbname=postgres user=postgres'

  - name: remote
    dsn: 'postgres://postgres:postgres@203.0.113.10:5432/postgres'
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/postgres.conf).

## Troubleshooting

To troubleshoot issues with the `postgres` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m postgres
  ```
