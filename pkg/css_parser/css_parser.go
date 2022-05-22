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
	rules  []*Rule
	tokens []tokenInfo
	cur    int
}

type Rule struct {
	Selector     string
	Declarations []Declaration
}

type Declaration struct {
	Property string
	Value    string
}

type tokenInfo struct {
	tokenType TokenType
	value     string
	start     int
	len       int
}

func newTokenInfo(tokenType TokenType, value string, start, len int) tokenInfo {
	return tokenInfo{tokenType, value, start, len}
}

func (t tokenInfo) String() string {
	return fmt.Sprintf("%s: %s[%d, %d]", t.tokenType, t.value, t.start, t.len)
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
	t.tokens = append(t.tokens, newTokenInfo(OpenCurlyBraceToken, "{", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexCloseCurlyBrace() {
	t.tokens = append(t.tokens, newTokenInfo(CloseCurlyBraceToken, "}", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexOpenSquareBracket() {
	t.tokens = append(t.tokens, newTokenInfo(OpenSquareBracketToken, "[", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexCloseSquareBracket() {
	t.tokens = append(t.tokens, newTokenInfo(CloseSquareBracketToken, "]", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexColon() {
	t.tokens = append(t.tokens, newTokenInfo(ColonToken, ":", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexSemicolon() {
	t.tokens = append(t.tokens, newTokenInfo(SemicolonToken, ";", t.cur, 1))
	t.advance()
}

func (t *tokenizer) lexComma() {
	t.tokens = append(t.tokens, newTokenInfo(CommaToken, ",", t.cur, 1))
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
			c == '-' || c == '_' || c == '.' || c == '%' {
			t.advance()
			ident += string(c)
		} else {
			break
		}
	}
	t.tokens = append(t.tokens, newTokenInfo(IdentifierToken, ident, t.cur-len(ident), len(ident)))
}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) Parse(css string) ([]*Rule, error) {
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

func (p *parser) peek() TokenType {
	if p.cur < len(p.tokens) {
		return p.tokens[p.cur].tokenType
	}
	return EOFToken
}

func (p *parser) matchRules() (*Rule, error) {
	rule := Rule{}
	if selector, err := p.matchSelector(); err != nil {
		return nil, fmt.Errorf("failed to match selector: %v", err)
	} else {
		rule.Selector = selector
	}
	if err := p.matchOpenCurlyBrace(); err != nil {
		return nil, err
	}
	if delcarations, err := p.matchDeclarations(); err != nil {
		return nil, fmt.Errorf("failed to match declarations: %v", err)
	} else {
		rule.Declarations = delcarations
	}
	if err := p.matchCloseCurlyBrace(); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (p *parser) matchDeclaration() (*Declaration, error) {
	delcaration := Declaration{}
	if property, err := p.matchProperty(); err != nil {
		return nil, err
	} else {
		delcaration.Property = property
	}
	if err := p.matchColon(); err != nil {
		return nil, err
	}
	if value, err := p.matchValue(); err != nil {
		return nil, err
	} else {
		delcaration.Value = value
	}
	if err := p.matchSemicolon(); err != nil {
		return nil, err
	}
	return &delcaration, nil
}

func (p *parser) matchDeclarations() ([]Declaration, error) {
	var declarations []Declaration
	for p.cur < len(p.tokens) {
		if declaration, err := p.matchDeclaration(); err != nil {
			return declarations, nil
		} else {
			declarations = append(declarations, *declaration)
		}
	}
	return declarations, nil
}

func (p *parser) matchProperty() (property string, err error) {
	return p.matchIdentifier()
}

func (p *parser) matchValue() (value string, err error) {
	value, err = p.matchIdentifier()
	if err != nil {
		return "", err
	}
	for p.peek() == IdentifierToken {
		if v, err := p.matchIdentifier(); err == nil {
			value += " " + v
		} else {
			return "", err
		}
	}
	return value, nil
}

func (p *parser) matchSelector() (string, error) {
	return p.matchIdentifier()
}

func (p *parser) matchIdentifier() (string, error) {
	if p.tokens[p.cur].tokenType == IdentifierToken {
		defer p.advance()
		return p.tokens[p.cur].value, nil
	}
	return "", fmt.Errorf("failed to match identifier: %v", p.tokens[p.cur])
}

func (p *parser) matchOpenCurlyBrace() error {
	if p.tokens[p.cur].tokenType == OpenCurlyBraceToken {
		p.advance()
		return nil
	}
	return fmt.Errorf("failed to match open curly brace: %v", p.tokens[p.cur])
}

func (p *parser) matchCloseCurlyBrace() error {
	if p.tokens[p.cur].tokenType == CloseCurlyBraceToken {
		p.advance()
		return nil
	}
	return fmt.Errorf("failed to match close curly brace: %v", p.tokens[p.cur])
}

func (p *parser) matchColon() error {
	if p.tokens[p.cur].tokenType == ColonToken {
		p.advance()
		return nil
	}
	return fmt.Errorf("failed to match colon: %v", p.tokens[p.cur])
}

func (p *parser) matchSemicolon() error {
	if p.tokens[p.cur].tokenType == SemicolonToken {
		p.advance()
		return nil
	}
	return fmt.Errorf("failed to match semicolon: %v", p.tokens[p.cur])
}
