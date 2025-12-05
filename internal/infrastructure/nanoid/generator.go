package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Generator struct {
	alphabet string
	length   int
}

// NewGenerator создает генератор с кастомным алфавитом (URL-friendly) и длиной.
func NewGenerator() *Generator {
	return &Generator{
		// Исключаем похожие символы (l, 1, I, O, 0), чтобы ID легко читались/печатались
		alphabet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		length:   8,
	}
}

func (g *Generator) Generate() string {
	id, _ := gonanoid.Generate(g.alphabet, g.length)
	return id
}
