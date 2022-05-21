package cssparser

import "fmt"

type TokenType uint32

const (
	ErrorToken TokenType = iota
	IdentifierToken
	OpenCurlyBraceToken
	CloseCurlyBraceToken
	OpenSquareBracketToken
	CloseSquareBracketToken
	ColonToken
	SemicolonToken
	CommaToken
	WhitespaceToken
	EOFToken
)

func (t TokenType) String() string {
	switch t {
	case ErrorToken:
		return "ErrorToken"
	case IdentifierToken:
		return "IdentifierToken"
	case OpenCurlyBraceToken:
		return "OpenCurlyBraceToken"
	case CloseCurlyBraceToken:
		return "CloseCurlyBraceToken"
	case OpenSquareBracketToken:
		return "OpenSquareBracketToken"
	case CloseSquareBracketToken:
		return "CloseSquareBracketToken"
	case ColonToken:
		return "ColonToken"
	case SemicolonToken:
		return "SemicolonToken"
	case CommaToken:
		return "CommaToken"
	case WhitespaceToken:
		return "WhitespaceToken"
	case EOFToken:
		return "EOFToken"
	default:
		return "UnknownToken"
	}
}

type parser struct {
	rules  []*rule
	tokens []tokenInfo
	cur    int
}

type rule struct {
	selector     string
	declarations []declaration
}

type declaration struct {
	property string
	value    string
}

type tokenInfo struct {
	tokenType TokenType
	value     string
}

func (t tokenInfo) String() string {
	return t.tokenType.String() + ": " + t.value
}

type tokenizer struct {
	source string
	tokens []tokenInfo
	cur    int
}

func NewTokenizer(source string) *tokenizer {
	return &tokenizer{source, nil, 0}
}

func (t *tokenizer) GetTokens() []tokenInfo {
	return t.tokens
}

func (t *tokenizer) advance() {
	t.cur++
}

func (t *tokenizer) next() rune {
	return rune(t.source[t.cur])
}

func (t *tokenizer) lex() {
	for t.cur < len(t.source) {
		c := t.next()
		switch {
		case c == '{':
			t.lexOpenCurlyBrace()
		case c == '}':
			t.lexCloseCurlyBrace()
		case c == '[':
			t.lexOpenSquareBracket()
		case c == ']':
			t.lexCloseSquareBracket()
		case c == ':':
			t.lexColon()
		case c == ';':
			t.lexSemicolon()
		case c == ',':
			t.lexComma()
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			t.lexWhitespace()
		default:
			t.lexIdentifier()
		}
	}
}

func (t *tokenizer) lexOpenCurlyBrace() {
	t.tokens = append(t.tokens, tokenInfo{OpenCurlyBraceToken, "{"})
	t.advance()
}

func (t *tokenizer) lexCloseCurlyBrace() {
	t.tokens = append(t.tokens, tokenInfo{CloseCurlyBraceToken, "}"})
	t.advance()
}

func (t *tokenizer) lexOpenSquareBracket() {
	t.tokens = append(t.tokens, tokenInfo{OpenSquareBracketToken, "["})
	t.advance()
}

func (t *tokenizer) lexCloseSquareBracket() {
	t.tokens = append(t.tokens, tokenInfo{CloseSquareBracketToken, "]"})
	t.advance()
}

func (t *tokenizer) lexColon() {
	t.tokens = append(t.tokens, tokenInfo{ColonToken, ":"})
	t.advance()
}

func (t *tokenizer) lexSemicolon() {
	t.tokens = append(t.tokens, tokenInfo{SemicolonToken, ";"})
	t.advance()
}

func (t *tokenizer) lexComma() {
	t.tokens = append(t.tokens, tokenInfo{CommaToken, ","})
	t.advance()
}

func (t *tokenizer) lexWhitespace() {
	// discard whitespace tokens
	t.advance()
}

func (t *tokenizer) lexIdentifier() {
	var ident string
	for t.cur < len(t.source) {
		c := t.next()
		if c >= 'a' && c <= 'z' ||
			c >= 'A' && c <= 'Z' ||
			c >= '0' && c <= '9' ||
			c == '-' || c == '_' || c == '.' {
			t.advance()
			ident += string(c)
		} else {
			break
		}
	}
	t.tokens = append(t.tokens, tokenInfo{IdentifierToken, ident})
}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) Parse(css string) ([]*rule, error) {
	tokenizer := NewTokenizer(css)
	tokenizer.lex()
	p.tokens = tokenizer.GetTokens()
	for p.cur < len(p.tokens) {
		if rule, err := p.matchRules(); err != nil {
			return nil, fmt.Errorf("failed to match rules: %v", err)
		} else {
			p.rules = append(p.rules, rule)
		}
	}
	return p.rules, nil
}

func (p *parser) advance() {
	p.cur++
}

func (p *parser) matchRules() (*rule, error) {
	rule := rule{}
	if match, selector := p.matchSelector(); !match {
		return nil, fmt.Errorf("failed to match selector")
	} else {
		rule.selector = selector
	}
	if !p.matchOpenCurlyBrace() {
		return nil, fmt.Errorf("failed to match open curly brace")
	}
	if delcarations, err := p.matchDeclarations(); err != nil {
		return nil, fmt.Errorf("failed to match declarations: %v", err)
	} else {
		rule.declarations = delcarations
	}
	if !p.matchCloseCurlyBrace() {
		return nil, fmt.Errorf("failed to match close curly brace")
	}
	return &rule, nil
}

func (p *parser) matchDeclaration() (*declaration, error) {
	delcaration := declaration{}
	if match, property := p.matchProperty(); !match {
		return nil, fmt.Errorf("failed to match property")
	} else {
		delcaration.property = property
	}
	if !p.matchColon() {
		return nil, fmt.Errorf("failed to match colon")
	}
	if match, value := p.matchValue(); !match {
		return nil, fmt.Errorf("failed to match value")
	} else {
		delcaration.value = value
	}
	if !p.matchSemicolon() {
		return nil, fmt.Errorf("failed to match semicolon")
	}
	return &delcaration, nil
}

func (p *parser) matchDeclarations() ([]declaration, error) {
	var declarations []declaration
	for p.cur < len(p.tokens) {
		if declaration, err := p.matchDeclaration(); err != nil {
			return declarations, nil
		} else {
			declarations = append(declarations, *declaration)
		}
	}
	return declarations, nil
}

func (p *parser) matchProperty() (bool, string) {
	return p.matchIdentifier()
}

func (p *parser) matchValue() (bool, string) {
	return p.matchIdentifier()
}

func (p *parser) matchSelector() (bool, string) {
	return p.matchIdentifier()
}

func (p *parser) matchIdentifier() (bool, string) {
	if p.tokens[p.cur].tokenType == IdentifierToken {
		defer p.advance()
		return true, p.tokens[p.cur].value
	}
	return false, ""
}

func (p *parser) matchOpenCurlyBrace() bool {
	if p.tokens[p.cur].tokenType == OpenCurlyBraceToken {
		p.advance()
		return true
	}
	return false
}

func (p *parser) matchCloseCurlyBrace() bool {
	if p.tokens[p.cur].tokenType == CloseCurlyBraceToken {
		p.advance()
		return true
	}
	return false
}

func (p *parser) matchColon() bool {
	if p.tokens[p.cur].tokenType == ColonToken {
		p.advance()
		return true
	}
	return false
}

func (p *parser) matchSemicolon() bool {
	if p.tokens[p.cur].tokenType == SemicolonToken {
		p.advance()
		return true
	}
	return false
}
