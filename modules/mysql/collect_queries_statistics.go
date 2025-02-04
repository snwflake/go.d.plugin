// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"context"
	"time"
)

// Table Schema:
// (MariaDB) https://mariadb.com/kb/en/information-schema-processlist-table/
// (MySql) https://dev.mysql.com/doc/refman/5.7/en/information-schema-processlist-table.html
const (
	queryInfoSchemaProcessList = `SELECT TIME,USER FROM INFORMATION_SCHEMA.PROCESSLIST 
WHERE Info IS NOT NULL AND Info NOT LIKE '%PROCESSLIST%' 
ORDER BY TIME`
)

func (m *MySQL) collectProcessListStatistics(collected map[string]int64) error {
	m.Debugf("executing query: '%s'", queryInfoSchemaProcessList)

	var (
		queryDuration              int64 // queryInfoSchemaProcessList execution time (not including row fetching)
		longestRunningQuerySeconds int64 // slowest query milliseconds in process list
		user                       string
	)

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout.Duration)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, queryInfoSchemaProcessList)
	if err != nil {
		return err
	}

	queryDuration = time.Since(start).Milliseconds()
	defer func() { _ = rows.Close() }()

	collected["process_list_queries_count_system"] = 0
	collected["process_list_queries_count_user"] = 0

	for rows.Next() {
		if err := rows.Scan(&longestRunningQuerySeconds, &user); err != nil {
			return err
		}
		// system user refers to non-client threads
		// event_scheduler is the thread used to monitor scheduled events
		// system user and event_scheduler threads are grouped as system/database threads
		// authenticated and unauthenticated user are grouped as users
		// please see USER section in
		// https://dev.mysql.com/doc/refman/8.0/en/information-schema-processlist-table.html
		switch user {
		case "system user", "event_scheduler":
			collected["process_list_queries_count_system"] += 1

		default:
			collected["process_list_queries_count_user"] += 1
		}
	}

	collected["process_list_fetch_query_duration"] = queryDuration
	collected["process_list_longest_query_duration"] = longestRunningQuerySeconds

	return nil
}
