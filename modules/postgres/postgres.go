// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"database/sql"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	module.Register("postgres", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Postgres {
	return &Postgres{
		Config: Config{
			Timeout: web.Duration{Duration: time.Second},
			DSN:     "postgres://postgres:postgres@127.0.0.1:5432/postgres",
		},
		charts:                 baseCharts.Copy(),
		recheckSettingsEvery:   time.Minute * 30,
		relistDatabaseEvery:    time.Minute,
		relistReplStandbyEvery: time.Minute,
		relistReplSlotEvery:    time.Minute,
	}
}

type Config struct {
	DSN     string       `yaml:"dsn"`
	Timeout web.Duration `yaml:"timeout"`
}

type Postgres struct {
	module.Base
	Config `yaml:",inline"`

	charts *module.Charts

	db *sql.DB

	superUser *bool
	pgVersion int

	maxConnections int64

	recheckSettingsTime    time.Time
	recheckSettingsEvery   time.Duration
	relistDatabaseTime     time.Time
	relistDatabaseEvery    time.Duration
	relistReplStandbyTime  time.Time
	relistReplStandbyEvery time.Duration
	relistReplSlotTime     time.Time
	relistReplSlotEvery    time.Duration

	databases       []string
	replStandbyApps []string
	replSlots       []string
}

func (p *Postgres) Init() bool {
	err := p.validateConfig()
	if err != nil {
		p.Errorf("config validation: %v", err)
		return false
	}

	return true
}

func (p *Postgres) Check() bool {
	return len(p.Collect()) > 0
}

func (p *Postgres) Charts() *module.Charts {
	return p.charts
}

func (p *Postgres) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (p *Postgres) Cleanup() {
	if p.db == nil {
		return
	}
	if err := p.db.Close(); err != nil {
		p.Warningf("cleanup: error on closing the Postgres database [%s]: %v", p.DSN, err)
	}
	p.db = nil
}
