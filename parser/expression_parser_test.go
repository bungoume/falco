package parser

import (
	"testing"

	"github.com/ysugimoto/falco/ast"
	"github.com/ysugimoto/falco/lexer"
)

func TestIfExpression(t *testing.T) {
	input := `// Subroutine
sub vcl_recv {
	// Leading comment
	set req.http.Foo = if (req.http.Host, "example.com", "foobar"); // Trailing comment
}`

	expect := &ast.VCL{
		Statements: []ast.Statement{
			&ast.SubroutineDeclaration{
				Meta: ast.New(T, 0, comments("// Subroutine")),
				Name: &ast.Ident{
					Meta:  ast.New(T, 0),
					Value: "vcl_recv",
				},
				Block: &ast.BlockStatement{
					Meta: ast.New(T, 1),
					Statements: []ast.Statement{
						&ast.SetStatement{
							Meta: ast.New(T, 1, comments("// Leading comment"), comments("// Trailing comment")),
							Ident: &ast.Ident{
								Meta:  ast.New(T, 1),
								Value: "req.http.Foo",
							},
							Operator: &ast.Operator{
								Meta:     ast.New(T, 1),
								Operator: "=",
							},
							Value: &ast.IfExpression{
								Meta: ast.New(T, 1),
								Condition: &ast.Ident{
									Meta:  ast.New(T, 1),
									Value: "req.http.Host",
								},
								Consequence: &ast.String{
									Meta:  ast.New(T, 1),
									Value: "example.com",
								},
								Alternative: &ast.String{
									Meta:  ast.New(T, 1),
									Value: "foobar",
								},
							},
						},
					},
				},
			},
		},
	}
	vcl, err := New(lexer.NewFromString(input)).ParseVCL()
	if err != nil {
		t.Errorf("%+v", err)
	}
	assert(t, vcl, expect)
}

func TestInfixIfExpression(t *testing.T) {
	input := `// Subroutine
sub vcl_recv {
	// Leading comment
	log {"foo bar"}
		if (req.http.Host, "example.com", "foobar")
		{"baz"}; // Trailing comment
}`

	expect := &ast.VCL{
		Statements: []ast.Statement{
			&ast.SubroutineDeclaration{
				Meta: ast.New(T, 0, comments("// Subroutine")),
				Name: &ast.Ident{
					Meta:  ast.New(T, 0),
					Value: "vcl_recv",
				},
				Block: &ast.BlockStatement{
					Meta: ast.New(T, 1),
					Statements: []ast.Statement{
						&ast.LogStatement{
							Meta: ast.New(T, 1, comments("// Leading comment"), comments("// Trailing comment")),
							Value: &ast.InfixExpression{
								Meta:     ast.New(T, 1),
								Operator: "+",
								Right: &ast.String{
									Meta:       ast.New(T, 1),
									Value:      "baz",
									LongString: true,
								},
								Left: &ast.InfixExpression{
									Meta:     ast.New(T, 1),
									Operator: "+",
									Left: &ast.String{
										Meta:       ast.New(T, 1),
										Value:      "foo bar",
										LongString: true,
									},
									Right: &ast.IfExpression{
										Meta: ast.New(T, 1),
										Condition: &ast.Ident{
											Meta:  ast.New(T, 1),
											Value: "req.http.Host",
										},
										Consequence: &ast.String{
											Meta:  ast.New(T, 1),
											Value: "example.com",
										},
										Alternative: &ast.String{
											Meta:  ast.New(T, 1),
											Value: "foobar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	vcl, err := New(lexer.NewFromString(input)).ParseVCL()
	if err != nil {
		t.Errorf("%+v", err)
	}
	assert(t, vcl, expect)
}

func TestFunctionCallExpression(t *testing.T) {
	t.Run("no argument", func(t *testing.T) {
		input := `// Subroutine
sub vcl_recv {
	set req.http.X-Trace-Id = uuid.version4();
}`
		expect := &ast.VCL{
			Statements: []ast.Statement{
				&ast.SubroutineDeclaration{
					Meta: ast.New(T, 0, comments("// Subroutine")),
					Name: &ast.Ident{
						Meta:  ast.New(T, 0),
						Value: "vcl_recv",
					},
					Block: &ast.BlockStatement{
						Meta: ast.New(T, 1),
						Statements: []ast.Statement{
							&ast.SetStatement{
								Meta: ast.New(T, 1),
								Ident: &ast.Ident{
									Meta:  ast.New(T, 1),
									Value: "req.http.X-Trace-Id",
								},
								Operator: &ast.Operator{
									Meta:     ast.New(T, 1),
									Operator: "=",
								},
								Value: &ast.FunctionCallExpression{
									Meta: ast.New(T, 1),
									Function: &ast.Ident{
										Meta:  ast.New(T, 1),
										Value: "uuid.version4",
									},
									Arguments: []ast.Expression{},
								},
							},
						},
					},
				},
			},
		}
		vcl, err := New(lexer.NewFromString(input)).ParseVCL()
		if err != nil {
			t.Errorf("%+v", err)
		}
		assert(t, vcl, expect)
	})

	t.Run("some arguments", func(t *testing.T) {
		input := `// Subroutine
sub vcl_recv {
	set req.http.X-Trace-Id = regsub(req.http.Host, "example.com", "");
}`
		expect := &ast.VCL{
			Statements: []ast.Statement{
				&ast.SubroutineDeclaration{
					Meta: ast.New(T, 0, comments("// Subroutine")),
					Name: &ast.Ident{
						Meta:  ast.New(T, 0),
						Value: "vcl_recv",
					},
					Block: &ast.BlockStatement{
						Meta: ast.New(T, 1),
						Statements: []ast.Statement{
							&ast.SetStatement{
								Meta: ast.New(T, 1),
								Ident: &ast.Ident{
									Meta:  ast.New(T, 1),
									Value: "req.http.X-Trace-Id",
								},
								Operator: &ast.Operator{
									Meta:     ast.New(T, 1),
									Operator: "=",
								},
								Value: &ast.FunctionCallExpression{
									Meta: ast.New(T, 1),
									Function: &ast.Ident{
										Meta:  ast.New(T, 1),
										Value: "regsub",
									},
									Arguments: []ast.Expression{
										&ast.Ident{
											Meta:  ast.New(T, 1),
											Value: "req.http.Host",
										},
										&ast.String{
											Meta:  ast.New(T, 1),
											Value: "example.com",
										},
										&ast.String{
											Meta:  ast.New(T, 1),
											Value: "",
										},
									},
								},
							},
						},
					},
				},
			},
		}
		vcl, err := New(lexer.NewFromString(input)).ParseVCL()
		if err != nil {
			t.Errorf("%+v", err)
		}
		assert(t, vcl, expect)
	})
}

func TestInfixExpression(t *testing.T) {
	t.Run("having infix comments in if condition", func(t *testing.T) {
		input := `
sub vcl_recv {
	// foo
	if (req.http.Foo && /* comment */ req.http.Bar) {
		// foo
		set req.http.Result = req.http.Foo req.http.Bar;
	}
	// end
}`
		expect := &ast.VCL{
			Statements: []ast.Statement{
				&ast.SubroutineDeclaration{
					Meta: ast.New(T, 0),
					Name: &ast.Ident{
						Meta:  ast.New(T, 0),
						Value: "vcl_recv",
					},
					Block: &ast.BlockStatement{
						Meta: ast.New(T, 1, ast.Comments{}, ast.Comments{}, comments("// end")),
						Statements: []ast.Statement{
							&ast.IfStatement{
								Meta: ast.New(T, 1, comments("// foo")),
								Condition: &ast.InfixExpression{
									Meta:     ast.New(T, 1),
									Operator: "&&",
									Left: &ast.Ident{
										Meta:  ast.New(T, 1),
										Value: "req.http.Foo",
									},
									Right: &ast.Ident{
										Meta:  ast.New(T, 1, comments("/* comment */")),
										Value: "req.http.Bar",
									},
								},
								Consequence: &ast.BlockStatement{
									Meta: ast.New(T, 2),
									Statements: []ast.Statement{
										&ast.SetStatement{
											Meta: ast.New(T, 2, comments("// foo")),
											Operator: &ast.Operator{
												Meta:     ast.New(T, 2),
												Operator: "=",
											},
											Ident: &ast.Ident{
												Meta:  ast.New(T, 2),
												Value: "req.http.Result",
											},
											Value: &ast.InfixExpression{
												Meta:     ast.New(T, 2),
												Operator: "+",
												Left: &ast.Ident{
													Meta:  ast.New(T, 2),
													Value: "req.http.Foo",
												},
												Right: &ast.Ident{
													Meta:  ast.New(T, 2),
													Value: "req.http.Bar",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		vcl, err := New(lexer.NewFromString(input)).ParseVCL()
		if err != nil {
			t.Errorf("%+v", err)
		}
		assert(t, vcl, expect)
	})
}
