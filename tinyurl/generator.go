package tinyurl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/unsubble/tiny-url/common"
	"github.com/unsubble/tiny-url/database"
)

var uniqueCode atomic.Int32

const (
	section1 uint64 = 0xFFFF00000
	section2 uint64 = 0xFFFF
)

type TinyURL common.TinyURL
type connChannel chan *sql.Conn
type ID int8

type Generator struct {
	ctx      *context.Context
	config   *config
	mu       sync.Mutex
	base     string
	lastCode string
	driver   database.Driver
}

type config struct {
	epochTime  uint64
	genId      ID
	dbId       ID
	uniqueCode ID
	dbName     string
	domain     string
}

type TimeoutError struct {
	msg string
}

func (gen *Generator) tryConnect(ch connChannel) {
	pool := gen.driver.GetPool()
	conn, err := pool.Conn(*gen.ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to create connection %s: %v\n", gen.config.dbName, err)
		ch <- nil
		return
	}
	ch <- conn
}

func (gen *Generator) getConnection(timeout time.Duration) (*sql.Conn, error) {
	ch := make(connChannel, 1)
	go func() {
		gen.tryConnect(ch)
		close(ch)
	}()

	select {
	case conn := <-ch:
		if conn == nil {
			return nil, NewTimeoutError("Database connection timed out")
		}
		return conn, nil
	case <-time.After(timeout):
		return nil, NewTimeoutError("Database connection timed out")
	}
}

func (gen *Generator) GenerateURL(targetUrl string) (*TinyURL, error) {
	conn, err := gen.getConnection(time.Second * 5)
	if err != nil {
		log.Default().Printf("[ERROR]: %v\n", err)
		return nil, err
	}
	defer conn.Close()

	gen.mu.Lock()
	shortCode := generateShortCode(gen.base, gen.config)
	for shortCode == gen.lastCode {
		shortCode = generateShortCode(gen.base, gen.config)
	}
	gen.lastCode = shortCode
	gen.mu.Unlock()

	generatedUrl := gen.concatUrl(shortCode)

	tinyUrl := &TinyURL{
		OriginalURL:  targetUrl,
		GeneratedURL: generatedUrl,
		GeneratedAt:  time.Now(),
	}

	gen.driver.AddUrl((*database.TinyURL)(tinyUrl))

	return tinyUrl, nil
}

func (gen *Generator) concatUrl(shortCode string) string {
	builder := strings.Builder{}
	if !strings.HasPrefix(gen.config.domain, "http://") &&
		!strings.HasPrefix(gen.config.domain, "https://") {
		builder.WriteString("https://")
	}
	builder.WriteString(gen.config.domain)
	if !strings.HasSuffix(gen.config.domain, "/") {
		builder.WriteRune('/')
	}
	builder.WriteString(shortCode)
	return builder.String()
}

func generateShortCode(base string, cfg *config) string {
	now := uint64(time.Now().UnixMilli())
	now -= cfg.epochTime

	intCode := now | (section1 & section2)
	intCode = intCode<<24 | uint64(cfg.genId)<<16 | uint64(cfg.dbId)<<8 | uint64(cfg.uniqueCode)

	length := uint64(len(base))

	builder := strings.Builder{}
	for intCode > 0 {
		letter := base[intCode%length]
		builder.WriteByte(letter)
		intCode /= length
	}

	return builder.String()
}

func NewGenerator(ctx *context.Context, cfg *common.Config, driver database.Driver, dbId ID, genId ID) *Generator {
	driverInfo := driver.GetInfo()
	innerConfig := &config{
		epochTime:  cfg.EpochTime,
		genId:      genId,
		dbId:       dbId,
		dbName:     fmt.Sprintf("%s\\%s", driverInfo.DriverName, driverInfo.DbName),
		uniqueCode: ID(uniqueCode.Add(1)),
		domain:     cfg.Domain,
	}
	return &Generator{
		ctx:      ctx,
		config:   innerConfig,
		mu:       sync.Mutex{},
		base:     cfg.BaseFormat,
		lastCode: "",
		driver:   driver,
	}
}

func NewTimeoutError(msg string) *TimeoutError {
	return &TimeoutError{
		msg: msg,
	}
}

func (t *TimeoutError) Error() string {
	return t.msg
}
