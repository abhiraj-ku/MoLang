package lexer

import "testing"

func TestTokenizeVariableDeclaration(t *testing.T) {
	tokens := Tokenize("yaha age = 20")

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{TokenKwVar, "yaha"},
		{TokenId, "age"},
		{TokenAssign, "="},
		{TokenNumber, "20"},
		{TokenEOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for index, want := range expected {
		got := tokens[index]
		if got.Type != want.tokenType || got.Value != want.value {
			t.Errorf("token %d: expected (%s, %q), got (%s, %q)", index, want.tokenType, want.value, got.Type, got.Value)
		}
	}
}

func TestTokenizeKeywordPrefixAsIdentifier(t *testing.T) {
	tokens := Tokenize("yahoo")

	if tokens[0].Type != TokenId || tokens[0].Value != "yahoo" {
		t.Fatalf("expected yahoo to be an identifier, got %#v", tokens[0])
	}
}
