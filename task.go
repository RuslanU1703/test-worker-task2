package workertask

import (
	"log"
	"task/netip"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

// Библиотека google/cel-go предоставляет удобный интерфейс для проверки Protobuf структуры на набор условий:

// https://github.com/google/cel-spec/blob/master/doc/intro.md
// https://github.com/google/cel-spec/blob/master/doc/langdef.md
// Однако часто возникает необходимость использовать Cel-GO для проверки обычных структур языка Go. Автор библиотеки дал подсказку как это можно реализовать: https://github.com/google/cel-go/issues/408

// Предлагаем вам:
// Реализовать ref.TypeProvider для следующей структуры:
type SecurityEvent struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	Tags      map[string]string      `json:"tags"`
	HostIP    netip.IP               `json:"host_ip"`
	Port      uint16                 `json:"port"`
	Flags     []string               `json:"flags"`
	Custom    map[string]interface{} `json:"custom"`
}

// my implementation
var program cel.Program

func init() {
	filter := `securityEvent.id.matches("[0-9-]") && "message" in securityEvent.tags`
	options := cel.Declarations(
		decls.NewVar("securityEvent", decls.NewMapType(decls.String, decls.Dyn)),
	)
	env, err := cel.NewEnv(options)
	if err != nil {
		log.Fatalln(err)
	}
	ast, iss := env.Compile(filter)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalln(err)
	}
	program = prg
}

// Написать юнит-тесты, в которых на вход будут подаваться заполненные структуры SecurityEvent и выражения Cel Go, на выходе - true если структура соответствует набору условий, false если не соответствует.
// Структура Custom содержит значения лишь перечисленных типов: int, string, bool, другими типами и уровнями вложенности принебречь.

// Результат отправить нам по почте, либо выложить на Github и прислать нам ссылку.
