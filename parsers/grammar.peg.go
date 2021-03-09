package parsers

import (
	"errors"
	"strings"
	"github.com/flyx/askew/data"
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulee
	ruleassignments
	rulebindings
	rulebinding
	ruleautovar
	ruletypedvar
	ruleisp
	ruleassignment
	rulebound
	ruleself
	ruledataset
	ruleprop
	rulestyle
	ruleclass
	ruleform
	ruleevent
	rulehtmlid
	rulejsid
	ruleexpr
	rulecommaless
	rulenumber
	ruleoperators
	rulestring
	ruleenclosed
	ruleparens
	rulebraces
	rulebrackets
	ruleinner
	ruleidentifier
	rulefields
	rulefsep
	rulefield
	rulename
	ruletype
	rulesname
	ruleqname
	rulearray
	rulemap
	rulechan
	rulefunc
	rulekeytype
	rulepointer
	rulecaptures
	rulecapture
	rulehandlername
	ruleeventid
	rulemappings
	rulemapping
	rulemappingname
	ruletags
	ruletag
	ruletagname
	ruletagarg
	rulefor
	ruleforVar
	rulehandlers
	rulehandler
	ruleparamname
	ruleparam
	rulecparams
	rulecparam
	rulevar
	ruleargs
	rulearg
	ruleimports
	ruleimport
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"e",
	"assignments",
	"bindings",
	"binding",
	"autovar",
	"typedvar",
	"isp",
	"assignment",
	"bound",
	"self",
	"dataset",
	"prop",
	"style",
	"class",
	"form",
	"event",
	"htmlid",
	"jsid",
	"expr",
	"commaless",
	"number",
	"operators",
	"string",
	"enclosed",
	"parens",
	"braces",
	"brackets",
	"inner",
	"identifier",
	"fields",
	"fsep",
	"field",
	"name",
	"type",
	"sname",
	"qname",
	"array",
	"map",
	"chan",
	"func",
	"keytype",
	"pointer",
	"captures",
	"capture",
	"handlername",
	"eventid",
	"mappings",
	"mapping",
	"mappingname",
	"tags",
	"tag",
	"tagname",
	"tagarg",
	"for",
	"forVar",
	"handlers",
	"handler",
	"paramname",
	"param",
	"cparams",
	"cparam",
	"var",
	"args",
	"arg",
	"imports",
	"import",
	"Action0",
	"PegText",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type GeneralParser struct {
	eventHandling              data.EventHandling
	expr, tagname, handlername string
	paramnames                 []string
	names                      []string
	keytype, valuetype         *data.ParamType
	fields                     []*data.Field
	bv                         data.BoundValue
	goVal                      data.GoValue
	paramMappings              map[string]data.BoundValue
	params                     []data.Param
	isVar                      bool
	err                        error

	assignments   []data.Assignment
	varMappings   []data.VariableMapping
	eventMappings []data.UnboundEventMapping
	handlers      []HandlerSpec
	cParams       []data.ComponentParam
	imports       map[string]string

	Buffer string
	buffer []rune
	rules  [108]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *GeneralParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *GeneralParser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *GeneralParser) Highlighter() {
	p.PrintSyntax()
}

func (p *GeneralParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.varMappings = append(p.varMappings,
				data.VariableMapping{Value: p.bv, Variable: p.goVal})
			p.goVal.Type = nil
			p.bv.IDs = nil

		case ruleAction1:

			p.goVal.Name = buffer[begin:end]

		case ruleAction2:

			p.goVal.Type = p.valuetype
			p.valuetype = nil

		case ruleAction3:

			p.assignments = append(p.assignments, data.Assignment{Expression: p.expr,
				Target: p.bv})
			p.bv.IDs = nil

		case ruleAction4:

			p.bv.Kind = data.BoundSelf

		case ruleAction5:

			p.bv.Kind = data.BoundDataset

		case ruleAction6:

			p.bv.Kind = data.BoundProperty

		case ruleAction7:

			p.bv.Kind = data.BoundStyle

		case ruleAction8:

			p.bv.Kind = data.BoundClass

		case ruleAction9:

			p.bv.Kind = data.BoundFormValue

		case ruleAction10:

			p.bv.Kind = data.BoundEventValue
			if len(p.bv.IDs) == 0 {
				p.bv.IDs = append(p.bv.IDs, "")
			}

		case ruleAction11:

			p.bv.IDs = append(p.bv.IDs, buffer[begin:end])

		case ruleAction12:

			p.bv.IDs = append(p.bv.IDs, buffer[begin:end])

		case ruleAction13:

			p.expr = buffer[begin:end]

		case ruleAction14:

			var expr *string
			if p.expr != "" {
				expr = new(string)
				*expr = p.expr
			}
			for _, name := range p.names {
				p.fields = append(p.fields, &data.Field{Name: name, Type: p.valuetype, DefaultValue: expr})
			}
			p.expr = ""
			p.valuetype = nil
			p.names = nil

		case ruleAction15:

			p.names = append(p.names, buffer[begin:end])

		case ruleAction16:

			switch name := buffer[begin:end]; name {
			case "int":
				p.valuetype = &data.ParamType{Kind: data.IntType}
			case "bool":
				p.valuetype = &data.ParamType{Kind: data.BoolType}
			case "string":
				p.valuetype = &data.ParamType{Kind: data.StringType}
			default:
				p.valuetype = &data.ParamType{Kind: data.NamedType, Name: name}
			}

		case ruleAction17:

			name := buffer[begin:end]
			if name == "js.Value" {
				p.valuetype = &data.ParamType{Kind: data.JSValueType}
			} else {
				p.valuetype = &data.ParamType{Kind: data.NamedType, Name: name}
			}

		case ruleAction18:

			p.valuetype = &data.ParamType{Kind: data.ArrayType, ValueType: p.valuetype}

		case ruleAction19:

			p.valuetype = &data.ParamType{Kind: data.MapType, KeyType: p.keytype, ValueType: p.valuetype}

		case ruleAction20:

			p.valuetype = &data.ParamType{Kind: data.ChanType, ValueType: p.valuetype}

		case ruleAction21:

			p.valuetype = &data.ParamType{Kind: data.FuncType, ValueType: p.valuetype,
				Params: p.params}
			p.params = nil

		case ruleAction22:

			p.keytype = p.valuetype

		case ruleAction23:

			p.valuetype = &data.ParamType{Kind: data.PointerType, ValueType: p.valuetype}

		case ruleAction24:

			p.eventMappings = append(p.eventMappings, data.UnboundEventMapping{
				Event: p.expr, Handler: p.handlername, ParamMappings: p.paramMappings,
				Handling: p.eventHandling})
			p.eventHandling = data.AutoPreventDefault
			p.expr = ""
			p.paramMappings = make(map[string]data.BoundValue)

		case ruleAction25:

			p.handlername = buffer[begin:end]

		case ruleAction26:

			p.expr = buffer[begin:end]

		case ruleAction27:

			if _, ok := p.paramMappings[p.tagname]; ok {
				p.err = errors.New("duplicate param: " + p.tagname)
				return
			}
			p.paramMappings[p.tagname] = p.bv
			p.bv.IDs = nil

		case ruleAction28:

			p.tagname = buffer[begin:end]

		case ruleAction29:

			switch p.tagname {
			case "preventDefault":
				if p.eventHandling != data.AutoPreventDefault {
					p.err = errors.New("duplicate preventDefault")
					return
				}
				switch len(p.names) {
				case 0:
					p.eventHandling = data.PreventDefault
				case 1:
					switch p.names[0] {
					case "true":
						p.eventHandling = data.PreventDefault
					case "false":
						p.eventHandling = data.DontPreventDefault
					case "ask":
						p.eventHandling = data.AskPreventDefault
					default:
						p.err = fmt.Errorf("unsupported value for preventDefault: %s", p.names[0])
						return
					}
				default:
					p.err = errors.New("too many parameters for preventDefault")
					return
				}
			default:
				p.err = errors.New("unknown tag: " + p.tagname)
				return
			}
			p.names = nil

		case ruleAction30:

			p.tagname = buffer[begin:end]

		case ruleAction31:

			p.names = append(p.names, buffer[begin:end])

		case ruleAction32:

			p.names = append(p.names, buffer[begin:end])

		case ruleAction33:

			p.handlers = append(p.handlers, HandlerSpec{
				Name: p.handlername, Params: p.params, Returns: p.valuetype})
			p.valuetype = nil
			p.params = nil

		case ruleAction34:

			p.paramnames = append(p.paramnames, buffer[begin:end])

		case ruleAction35:

			name := p.paramnames[len(p.paramnames)-1]
			p.paramnames = p.paramnames[:len(p.paramnames)-1]
			for _, para := range p.params {
				if para.Name == name {
					p.err = errors.New("duplicate param name: " + para.Name)
					return
				}
			}

			p.params = append(p.params, data.Param{Name: name, Type: p.valuetype})
			p.valuetype = nil

		case ruleAction36:

			p.cParams = append(p.cParams, data.ComponentParam{
				Name: p.tagname, Type: *p.valuetype, IsVar: p.isVar})
			p.valuetype = nil
			p.isVar = false

		case ruleAction37:

			p.isVar = true

		case ruleAction38:

			p.names = append(p.names, p.expr)

		case ruleAction39:

			path := buffer[begin:end]
			if p.tagname == "" {
				lastDot := strings.LastIndexByte(path, '/')
				if lastDot == -1 {
					p.tagname = path
				} else {
					p.tagname = path[lastDot+1:]
				}
			}
			if _, ok := p.imports[p.tagname]; ok {
				p.err = errors.New("duplicate import name: " + p.tagname)
				return
			}
			p.imports[p.tagname] = path
			p.tagname = ""

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *GeneralParser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 e <- <(assignments / bindings / captures / fields / for / handlers / cparams / args / imports)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleassignments]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulebindings]() {
						goto l4
					}
					goto l2
				l4:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulecaptures]() {
						goto l5
					}
					goto l2
				l5:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulefields]() {
						goto l6
					}
					goto l2
				l6:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulefor]() {
						goto l7
					}
					goto l2
				l7:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulehandlers]() {
						goto l8
					}
					goto l2
				l8:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulecparams]() {
						goto l9
					}
					goto l2
				l9:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleargs]() {
						goto l10
					}
					goto l2
				l10:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleimports]() {
						goto l0
					}
				}
			l2:
				depth--
				add(rulee, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 assignments <- <(isp* assignment isp* ((',' / ';') isp* assignment isp*)* !.)> */
		func() bool {
			position11, tokenIndex11, depth11 := position, tokenIndex, depth
			{
				position12 := position
				depth++
			l13:
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l14
					}
					goto l13
				l14:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
				}
				if !_rules[ruleassignment]() {
					goto l11
				}
			l15:
				{
					position16, tokenIndex16, depth16 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l16
					}
					goto l15
				l16:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
				}
			l17:
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					{
						position19, tokenIndex19, depth19 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l20
						}
						position++
						goto l19
					l20:
						position, tokenIndex, depth = position19, tokenIndex19, depth19
						if buffer[position] != rune(';') {
							goto l18
						}
						position++
					}
				l19:
				l21:
					{
						position22, tokenIndex22, depth22 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l22
						}
						goto l21
					l22:
						position, tokenIndex, depth = position22, tokenIndex22, depth22
					}
					if !_rules[ruleassignment]() {
						goto l18
					}
				l23:
					{
						position24, tokenIndex24, depth24 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l24
						}
						goto l23
					l24:
						position, tokenIndex, depth = position24, tokenIndex24, depth24
					}
					goto l17
				l18:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
				}
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if !matchDot() {
						goto l25
					}
					goto l11
				l25:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
				}
				depth--
				add(ruleassignments, position12)
			}
			return true
		l11:
			position, tokenIndex, depth = position11, tokenIndex11, depth11
			return false
		},
		/* 2 bindings <- <(isp* binding isp* ((',' / ';') isp* binding isp*)* !.)> */
		func() bool {
			position26, tokenIndex26, depth26 := position, tokenIndex, depth
			{
				position27 := position
				depth++
			l28:
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l29
					}
					goto l28
				l29:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
				}
				if !_rules[rulebinding]() {
					goto l26
				}
			l30:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
				}
			l32:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					{
						position34, tokenIndex34, depth34 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l35
						}
						position++
						goto l34
					l35:
						position, tokenIndex, depth = position34, tokenIndex34, depth34
						if buffer[position] != rune(';') {
							goto l33
						}
						position++
					}
				l34:
				l36:
					{
						position37, tokenIndex37, depth37 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l37
						}
						goto l36
					l37:
						position, tokenIndex, depth = position37, tokenIndex37, depth37
					}
					if !_rules[rulebinding]() {
						goto l33
					}
				l38:
					{
						position39, tokenIndex39, depth39 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l39
						}
						goto l38
					l39:
						position, tokenIndex, depth = position39, tokenIndex39, depth39
					}
					goto l32
				l33:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
				}
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if !matchDot() {
						goto l40
					}
					goto l26
				l40:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
				}
				depth--
				add(rulebindings, position27)
			}
			return true
		l26:
			position, tokenIndex, depth = position26, tokenIndex26, depth26
			return false
		},
		/* 3 binding <- <(bound isp* ':' isp* (autovar / typedvar) Action0)> */
		func() bool {
			position41, tokenIndex41, depth41 := position, tokenIndex, depth
			{
				position42 := position
				depth++
				if !_rules[rulebound]() {
					goto l41
				}
			l43:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l44
					}
					goto l43
				l44:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
				}
				if buffer[position] != rune(':') {
					goto l41
				}
				position++
			l45:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l46
					}
					goto l45
				l46:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
				}
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if !_rules[ruleautovar]() {
						goto l48
					}
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if !_rules[ruletypedvar]() {
						goto l41
					}
				}
			l47:
				if !_rules[ruleAction0]() {
					goto l41
				}
				depth--
				add(rulebinding, position42)
			}
			return true
		l41:
			position, tokenIndex, depth = position41, tokenIndex41, depth41
			return false
		},
		/* 4 autovar <- <(<identifier> Action1)> */
		func() bool {
			position49, tokenIndex49, depth49 := position, tokenIndex, depth
			{
				position50 := position
				depth++
				{
					position51 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l49
					}
					depth--
					add(rulePegText, position51)
				}
				if !_rules[ruleAction1]() {
					goto l49
				}
				depth--
				add(ruleautovar, position50)
			}
			return true
		l49:
			position, tokenIndex, depth = position49, tokenIndex49, depth49
			return false
		},
		/* 5 typedvar <- <('(' isp* autovar isp+ type isp* ')' Action2)> */
		func() bool {
			position52, tokenIndex52, depth52 := position, tokenIndex, depth
			{
				position53 := position
				depth++
				if buffer[position] != rune('(') {
					goto l52
				}
				position++
			l54:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
				if !_rules[ruleautovar]() {
					goto l52
				}
				if !_rules[ruleisp]() {
					goto l52
				}
			l56:
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l57
					}
					goto l56
				l57:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
				}
				if !_rules[ruletype]() {
					goto l52
				}
			l58:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l59
					}
					goto l58
				l59:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
				}
				if buffer[position] != rune(')') {
					goto l52
				}
				position++
				if !_rules[ruleAction2]() {
					goto l52
				}
				depth--
				add(ruletypedvar, position53)
			}
			return true
		l52:
			position, tokenIndex, depth = position52, tokenIndex52, depth52
			return false
		},
		/* 6 isp <- <(' ' / '\t')> */
		func() bool {
			position60, tokenIndex60, depth60 := position, tokenIndex, depth
			{
				position61 := position
				depth++
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('\t') {
						goto l60
					}
					position++
				}
			l62:
				depth--
				add(ruleisp, position61)
			}
			return true
		l60:
			position, tokenIndex, depth = position60, tokenIndex60, depth60
			return false
		},
		/* 7 assignment <- <(isp* bound isp* '=' isp* expr Action3)> */
		func() bool {
			position64, tokenIndex64, depth64 := position, tokenIndex, depth
			{
				position65 := position
				depth++
			l66:
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l67
					}
					goto l66
				l67:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
				}
				if !_rules[rulebound]() {
					goto l64
				}
			l68:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l69
					}
					goto l68
				l69:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
				}
				if buffer[position] != rune('=') {
					goto l64
				}
				position++
			l70:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l71
					}
					goto l70
				l71:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
				}
				if !_rules[ruleexpr]() {
					goto l64
				}
				if !_rules[ruleAction3]() {
					goto l64
				}
				depth--
				add(ruleassignment, position65)
			}
			return true
		l64:
			position, tokenIndex, depth = position64, tokenIndex64, depth64
			return false
		},
		/* 8 bound <- <(self / ((&('E' | 'e') event) | (&('F' | 'f') form) | (&('C' | 'c') class) | (&('S' | 's') style) | (&('P' | 'p') prop) | (&('D' | 'd') dataset)))> */
		func() bool {
			position72, tokenIndex72, depth72 := position, tokenIndex, depth
			{
				position73 := position
				depth++
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if !_rules[ruleself]() {
						goto l75
					}
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					{
						switch buffer[position] {
						case 'E', 'e':
							if !_rules[ruleevent]() {
								goto l72
							}
							break
						case 'F', 'f':
							if !_rules[ruleform]() {
								goto l72
							}
							break
						case 'C', 'c':
							if !_rules[ruleclass]() {
								goto l72
							}
							break
						case 'S', 's':
							if !_rules[rulestyle]() {
								goto l72
							}
							break
						case 'P', 'p':
							if !_rules[ruleprop]() {
								goto l72
							}
							break
						default:
							if !_rules[ruledataset]() {
								goto l72
							}
							break
						}
					}

				}
			l74:
				depth--
				add(rulebound, position73)
			}
			return true
		l72:
			position, tokenIndex, depth = position72, tokenIndex72, depth72
			return false
		},
		/* 9 self <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('f' / 'F') isp* '(' isp* ')' Action4)> */
		func() bool {
			position77, tokenIndex77, depth77 := position, tokenIndex, depth
			{
				position78 := position
				depth++
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
					if buffer[position] != rune('S') {
						goto l77
					}
					position++
				}
			l79:
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('E') {
						goto l77
					}
					position++
				}
			l81:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('L') {
						goto l77
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('F') {
						goto l77
					}
					position++
				}
			l85:
			l87:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l88
					}
					goto l87
				l88:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
				}
				if buffer[position] != rune('(') {
					goto l77
				}
				position++
			l89:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l90
					}
					goto l89
				l90:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
				}
				if buffer[position] != rune(')') {
					goto l77
				}
				position++
				if !_rules[ruleAction4]() {
					goto l77
				}
				depth--
				add(ruleself, position78)
			}
			return true
		l77:
			position, tokenIndex, depth = position77, tokenIndex77, depth77
			return false
		},
		/* 10 dataset <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('s' / 'S') ('e' / 'E') ('t' / 'T') isp* '(' isp* htmlid isp* ')' Action5)> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('D') {
						goto l91
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('T') {
						goto l91
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('S') {
						goto l91
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('E') {
						goto l91
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('T') {
						goto l91
					}
					position++
				}
			l105:
			l107:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l108
					}
					goto l107
				l108:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
				}
				if buffer[position] != rune('(') {
					goto l91
				}
				position++
			l109:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l110
					}
					goto l109
				l110:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
				}
				if !_rules[rulehtmlid]() {
					goto l91
				}
			l111:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l112
					}
					goto l111
				l112:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
				}
				if buffer[position] != rune(')') {
					goto l91
				}
				position++
				if !_rules[ruleAction5]() {
					goto l91
				}
				depth--
				add(ruledataset, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 11 prop <- <(('p' / 'P') ('r' / 'R') ('o' / 'O') ('p' / 'P') isp* '(' isp* htmlid isp* ')' Action6)> */
		func() bool {
			position113, tokenIndex113, depth113 := position, tokenIndex, depth
			{
				position114 := position
				depth++
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('P') {
						goto l113
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('R') {
						goto l113
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('O') {
						goto l113
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('P') {
						goto l113
					}
					position++
				}
			l121:
			l123:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l124
					}
					goto l123
				l124:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
				}
				if buffer[position] != rune('(') {
					goto l113
				}
				position++
			l125:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l126
					}
					goto l125
				l126:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
				}
				if !_rules[rulehtmlid]() {
					goto l113
				}
			l127:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l128
					}
					goto l127
				l128:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
				}
				if buffer[position] != rune(')') {
					goto l113
				}
				position++
				if !_rules[ruleAction6]() {
					goto l113
				}
				depth--
				add(ruleprop, position114)
			}
			return true
		l113:
			position, tokenIndex, depth = position113, tokenIndex113, depth113
			return false
		},
		/* 12 style <- <(('s' / 'S') ('t' / 'T') ('y' / 'Y') ('l' / 'L') ('e' / 'E') isp* '(' isp* htmlid isp* ')' Action7)> */
		func() bool {
			position129, tokenIndex129, depth129 := position, tokenIndex, depth
			{
				position130 := position
				depth++
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('S') {
						goto l129
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('T') {
						goto l129
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('Y') {
						goto l129
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('L') {
						goto l129
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('E') {
						goto l129
					}
					position++
				}
			l139:
			l141:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l142
					}
					goto l141
				l142:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
				}
				if buffer[position] != rune('(') {
					goto l129
				}
				position++
			l143:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l144
					}
					goto l143
				l144:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
				}
				if !_rules[rulehtmlid]() {
					goto l129
				}
			l145:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l146
					}
					goto l145
				l146:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
				}
				if buffer[position] != rune(')') {
					goto l129
				}
				position++
				if !_rules[ruleAction7]() {
					goto l129
				}
				depth--
				add(rulestyle, position130)
			}
			return true
		l129:
			position, tokenIndex, depth = position129, tokenIndex129, depth129
			return false
		},
		/* 13 class <- <(('c' / 'C') ('l' / 'L') ('a' / 'A') ('s' / 'S') ('s' / 'S') isp* '(' isp* htmlid isp* (',' isp* htmlid isp*)* ')' Action8)> */
		func() bool {
			position147, tokenIndex147, depth147 := position, tokenIndex, depth
			{
				position148 := position
				depth++
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('C') {
						goto l147
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('L') {
						goto l147
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('A') {
						goto l147
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('S') {
						goto l147
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('S') {
						goto l147
					}
					position++
				}
			l157:
			l159:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l160
					}
					goto l159
				l160:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
				}
				if buffer[position] != rune('(') {
					goto l147
				}
				position++
			l161:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l162
					}
					goto l161
				l162:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
				}
				if !_rules[rulehtmlid]() {
					goto l147
				}
			l163:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l164
					}
					goto l163
				l164:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
				}
			l165:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l166
					}
					position++
				l167:
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l168
						}
						goto l167
					l168:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
					}
					if !_rules[rulehtmlid]() {
						goto l166
					}
				l169:
					{
						position170, tokenIndex170, depth170 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l170
						}
						goto l169
					l170:
						position, tokenIndex, depth = position170, tokenIndex170, depth170
					}
					goto l165
				l166:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
				}
				if buffer[position] != rune(')') {
					goto l147
				}
				position++
				if !_rules[ruleAction8]() {
					goto l147
				}
				depth--
				add(ruleclass, position148)
			}
			return true
		l147:
			position, tokenIndex, depth = position147, tokenIndex147, depth147
			return false
		},
		/* 14 form <- <(('f' / 'F') ('o' / 'O') ('r' / 'R') ('m' / 'M') isp* '(' isp* htmlid isp* ')' Action9)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('F') {
						goto l171
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('O') {
						goto l171
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('R') {
						goto l171
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('M') {
						goto l171
					}
					position++
				}
			l179:
			l181:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l182
					}
					goto l181
				l182:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
				}
				if buffer[position] != rune('(') {
					goto l171
				}
				position++
			l183:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l184
					}
					goto l183
				l184:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
				}
				if !_rules[rulehtmlid]() {
					goto l171
				}
			l185:
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l186
					}
					goto l185
				l186:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
				}
				if buffer[position] != rune(')') {
					goto l171
				}
				position++
				if !_rules[ruleAction9]() {
					goto l171
				}
				depth--
				add(ruleform, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 15 event <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('n' / 'N') ('t' / 'T') isp* '(' isp* jsid? isp* ')' Action10)> */
		func() bool {
			position187, tokenIndex187, depth187 := position, tokenIndex, depth
			{
				position188 := position
				depth++
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('E') {
						goto l187
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('V') {
						goto l187
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('E') {
						goto l187
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('N') {
						goto l187
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('T') {
						goto l187
					}
					position++
				}
			l197:
			l199:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
				}
				if buffer[position] != rune('(') {
					goto l187
				}
				position++
			l201:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l202
					}
					goto l201
				l202:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
				}
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if !_rules[rulejsid]() {
						goto l203
					}
					goto l204
				l203:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
				}
			l204:
			l205:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l206
					}
					goto l205
				l206:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
				}
				if buffer[position] != rune(')') {
					goto l187
				}
				position++
				if !_rules[ruleAction10]() {
					goto l187
				}
				depth--
				add(ruleevent, position188)
			}
			return true
		l187:
			position, tokenIndex, depth = position187, tokenIndex187, depth187
			return false
		},
		/* 16 htmlid <- <(<((&('-') '-') | (&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> Action11)> */
		func() bool {
			position207, tokenIndex207, depth207 := position, tokenIndex, depth
			{
				position208 := position
				depth++
				{
					position209 := position
					depth++
					{
						switch buffer[position] {
						case '-':
							if buffer[position] != rune('-') {
								goto l207
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l207
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l207
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l207
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l207
							}
							position++
							break
						}
					}

				l210:
					{
						position211, tokenIndex211, depth211 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '-':
								if buffer[position] != rune('-') {
									goto l211
								}
								position++
								break
							case '_':
								if buffer[position] != rune('_') {
									goto l211
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l211
								}
								position++
								break
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l211
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l211
								}
								position++
								break
							}
						}

						goto l210
					l211:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
					}
					depth--
					add(rulePegText, position209)
				}
				if !_rules[ruleAction11]() {
					goto l207
				}
				depth--
				add(rulehtmlid, position208)
			}
			return true
		l207:
			position, tokenIndex, depth = position207, tokenIndex207, depth207
			return false
		},
		/* 17 jsid <- <(<(((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> Action12)> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				{
					position216 := position
					depth++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l214
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l214
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l214
							}
							position++
							break
						}
					}

				l218:
					{
						position219, tokenIndex219, depth219 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l219
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l219
								}
								position++
								break
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l219
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l219
								}
								position++
								break
							}
						}

						goto l218
					l219:
						position, tokenIndex, depth = position219, tokenIndex219, depth219
					}
					depth--
					add(rulePegText, position216)
				}
				if !_rules[ruleAction12]() {
					goto l214
				}
				depth--
				add(rulejsid, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
			return false
		},
		/* 18 expr <- <(<((&('\t' | ' ') isp+) | (&('(' | '[' | '{') enclosed) | (&('!' | '"' | '&' | '*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | '<' | '=' | '>' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '|') commaless))+> Action13)> */
		func() bool {
			position221, tokenIndex221, depth221 := position, tokenIndex, depth
			{
				position222 := position
				depth++
				{
					position223 := position
					depth++
					{
						switch buffer[position] {
						case '\t', ' ':
							if !_rules[ruleisp]() {
								goto l221
							}
						l227:
							{
								position228, tokenIndex228, depth228 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l228
								}
								goto l227
							l228:
								position, tokenIndex, depth = position228, tokenIndex228, depth228
							}
							break
						case '(', '[', '{':
							if !_rules[ruleenclosed]() {
								goto l221
							}
							break
						default:
							if !_rules[rulecommaless]() {
								goto l221
							}
							break
						}
					}

				l224:
					{
						position225, tokenIndex225, depth225 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '\t', ' ':
								if !_rules[ruleisp]() {
									goto l225
								}
							l230:
								{
									position231, tokenIndex231, depth231 := position, tokenIndex, depth
									if !_rules[ruleisp]() {
										goto l231
									}
									goto l230
								l231:
									position, tokenIndex, depth = position231, tokenIndex231, depth231
								}
								break
							case '(', '[', '{':
								if !_rules[ruleenclosed]() {
									goto l225
								}
								break
							default:
								if !_rules[rulecommaless]() {
									goto l225
								}
								break
							}
						}

						goto l224
					l225:
						position, tokenIndex, depth = position225, tokenIndex225, depth225
					}
					depth--
					add(rulePegText, position223)
				}
				if !_rules[ruleAction13]() {
					goto l221
				}
				depth--
				add(ruleexpr, position222)
			}
			return true
		l221:
			position, tokenIndex, depth = position221, tokenIndex221, depth221
			return false
		},
		/* 19 commaless <- <((((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+ '.' ((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+) / ((&('"' | '`') string) | (&('!' | '&' | '*' | '+' | '-' | '.' | '/' | ':' | '<' | '=' | '>' | '^' | '|') operators) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') number) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') identifier)))> */
		func() bool {
			position232, tokenIndex232, depth232 := position, tokenIndex, depth
			{
				position233 := position
				depth++
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l235
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l235
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l235
							}
							position++
							break
						}
					}

				l236:
					{
						position237, tokenIndex237, depth237 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l237
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l237
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l237
								}
								position++
								break
							}
						}

						goto l236
					l237:
						position, tokenIndex, depth = position237, tokenIndex237, depth237
					}
					if buffer[position] != rune('.') {
						goto l235
					}
					position++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l235
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l235
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l235
							}
							position++
							break
						}
					}

				l240:
					{
						position241, tokenIndex241, depth241 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l241
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l241
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l241
								}
								position++
								break
							}
						}

						goto l240
					l241:
						position, tokenIndex, depth = position241, tokenIndex241, depth241
					}
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					{
						switch buffer[position] {
						case '"', '`':
							if !_rules[rulestring]() {
								goto l232
							}
							break
						case '!', '&', '*', '+', '-', '.', '/', ':', '<', '=', '>', '^', '|':
							if !_rules[ruleoperators]() {
								goto l232
							}
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[rulenumber]() {
								goto l232
							}
							break
						default:
							if !_rules[ruleidentifier]() {
								goto l232
							}
							break
						}
					}

				}
			l234:
				depth--
				add(rulecommaless, position233)
			}
			return true
		l232:
			position, tokenIndex, depth = position232, tokenIndex232, depth232
			return false
		},
		/* 20 number <- <[0-9]+> */
		func() bool {
			position245, tokenIndex245, depth245 := position, tokenIndex, depth
			{
				position246 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l245
				}
				position++
			l247:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
				}
				depth--
				add(rulenumber, position246)
			}
			return true
		l245:
			position, tokenIndex, depth = position245, tokenIndex245, depth245
			return false
		},
		/* 21 operators <- <((&('>') '>') | (&('<') '<') | (&('!') '!') | (&('.') '.') | (&('=') '=') | (&(':') ':') | (&('^') '^') | (&('&') '&') | (&('|') '|') | (&('/') '/') | (&('*') '*') | (&('-') '-') | (&('+') '+'))+> */
		func() bool {
			position249, tokenIndex249, depth249 := position, tokenIndex, depth
			{
				position250 := position
				depth++
				{
					switch buffer[position] {
					case '>':
						if buffer[position] != rune('>') {
							goto l249
						}
						position++
						break
					case '<':
						if buffer[position] != rune('<') {
							goto l249
						}
						position++
						break
					case '!':
						if buffer[position] != rune('!') {
							goto l249
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l249
						}
						position++
						break
					case '=':
						if buffer[position] != rune('=') {
							goto l249
						}
						position++
						break
					case ':':
						if buffer[position] != rune(':') {
							goto l249
						}
						position++
						break
					case '^':
						if buffer[position] != rune('^') {
							goto l249
						}
						position++
						break
					case '&':
						if buffer[position] != rune('&') {
							goto l249
						}
						position++
						break
					case '|':
						if buffer[position] != rune('|') {
							goto l249
						}
						position++
						break
					case '/':
						if buffer[position] != rune('/') {
							goto l249
						}
						position++
						break
					case '*':
						if buffer[position] != rune('*') {
							goto l249
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l249
						}
						position++
						break
					default:
						if buffer[position] != rune('+') {
							goto l249
						}
						position++
						break
					}
				}

			l251:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '>':
							if buffer[position] != rune('>') {
								goto l252
							}
							position++
							break
						case '<':
							if buffer[position] != rune('<') {
								goto l252
							}
							position++
							break
						case '!':
							if buffer[position] != rune('!') {
								goto l252
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l252
							}
							position++
							break
						case '=':
							if buffer[position] != rune('=') {
								goto l252
							}
							position++
							break
						case ':':
							if buffer[position] != rune(':') {
								goto l252
							}
							position++
							break
						case '^':
							if buffer[position] != rune('^') {
								goto l252
							}
							position++
							break
						case '&':
							if buffer[position] != rune('&') {
								goto l252
							}
							position++
							break
						case '|':
							if buffer[position] != rune('|') {
								goto l252
							}
							position++
							break
						case '/':
							if buffer[position] != rune('/') {
								goto l252
							}
							position++
							break
						case '*':
							if buffer[position] != rune('*') {
								goto l252
							}
							position++
							break
						case '-':
							if buffer[position] != rune('-') {
								goto l252
							}
							position++
							break
						default:
							if buffer[position] != rune('+') {
								goto l252
							}
							position++
							break
						}
					}

					goto l251
				l252:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
				}
				depth--
				add(ruleoperators, position250)
			}
			return true
		l249:
			position, tokenIndex, depth = position249, tokenIndex249, depth249
			return false
		},
		/* 22 string <- <(('`' (!'`' .)* '`') / ('"' ((!'"' .) / ('\\' '"'))* '"'))> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('`') {
						goto l258
					}
					position++
				l259:
					{
						position260, tokenIndex260, depth260 := position, tokenIndex, depth
						{
							position261, tokenIndex261, depth261 := position, tokenIndex, depth
							if buffer[position] != rune('`') {
								goto l261
							}
							position++
							goto l260
						l261:
							position, tokenIndex, depth = position261, tokenIndex261, depth261
						}
						if !matchDot() {
							goto l260
						}
						goto l259
					l260:
						position, tokenIndex, depth = position260, tokenIndex260, depth260
					}
					if buffer[position] != rune('`') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('"') {
						goto l255
					}
					position++
				l262:
					{
						position263, tokenIndex263, depth263 := position, tokenIndex, depth
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							{
								position266, tokenIndex266, depth266 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l266
								}
								position++
								goto l265
							l266:
								position, tokenIndex, depth = position266, tokenIndex266, depth266
							}
							if !matchDot() {
								goto l265
							}
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('\\') {
								goto l263
							}
							position++
							if buffer[position] != rune('"') {
								goto l263
							}
							position++
						}
					l264:
						goto l262
					l263:
						position, tokenIndex, depth = position263, tokenIndex263, depth263
					}
					if buffer[position] != rune('"') {
						goto l255
					}
					position++
				}
			l257:
				depth--
				add(rulestring, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 23 enclosed <- <((&('[') brackets) | (&('{') braces) | (&('(') parens))> */
		func() bool {
			position267, tokenIndex267, depth267 := position, tokenIndex, depth
			{
				position268 := position
				depth++
				{
					switch buffer[position] {
					case '[':
						if !_rules[rulebrackets]() {
							goto l267
						}
						break
					case '{':
						if !_rules[rulebraces]() {
							goto l267
						}
						break
					default:
						if !_rules[ruleparens]() {
							goto l267
						}
						break
					}
				}

				depth--
				add(ruleenclosed, position268)
			}
			return true
		l267:
			position, tokenIndex, depth = position267, tokenIndex267, depth267
			return false
		},
		/* 24 parens <- <('(' inner ')')> */
		func() bool {
			position270, tokenIndex270, depth270 := position, tokenIndex, depth
			{
				position271 := position
				depth++
				if buffer[position] != rune('(') {
					goto l270
				}
				position++
				if !_rules[ruleinner]() {
					goto l270
				}
				if buffer[position] != rune(')') {
					goto l270
				}
				position++
				depth--
				add(ruleparens, position271)
			}
			return true
		l270:
			position, tokenIndex, depth = position270, tokenIndex270, depth270
			return false
		},
		/* 25 braces <- <('{' inner '}')> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				if buffer[position] != rune('{') {
					goto l272
				}
				position++
				if !_rules[ruleinner]() {
					goto l272
				}
				if buffer[position] != rune('}') {
					goto l272
				}
				position++
				depth--
				add(rulebraces, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 26 brackets <- <('[' inner ']')> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				if buffer[position] != rune('[') {
					goto l274
				}
				position++
				if !_rules[ruleinner]() {
					goto l274
				}
				if buffer[position] != rune(']') {
					goto l274
				}
				position++
				depth--
				add(rulebrackets, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
			return false
		},
		/* 27 inner <- <((&('\t' | ' ') isp+) | (&(',') ',') | (&('(' | '[' | '{') enclosed) | (&('!' | '"' | '&' | '*' | '+' | '-' | '.' | '/' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | '<' | '=' | '>' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '^' | '_' | '`' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | '|') commaless))*> */
		func() bool {
			{
				position277 := position
				depth++
			l278:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\t', ' ':
							if !_rules[ruleisp]() {
								goto l279
							}
						l281:
							{
								position282, tokenIndex282, depth282 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l282
								}
								goto l281
							l282:
								position, tokenIndex, depth = position282, tokenIndex282, depth282
							}
							break
						case ',':
							if buffer[position] != rune(',') {
								goto l279
							}
							position++
							break
						case '(', '[', '{':
							if !_rules[ruleenclosed]() {
								goto l279
							}
							break
						default:
							if !_rules[rulecommaless]() {
								goto l279
							}
							break
						}
					}

					goto l278
				l279:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
				}
				depth--
				add(ruleinner, position277)
			}
			return true
		},
		/* 28 identifier <- <(((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') ([0-9] / [0-9])) | (&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> */
		func() bool {
			position283, tokenIndex283, depth283 := position, tokenIndex, depth
			{
				position284 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l283
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l283
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l283
						}
						position++
						break
					}
				}

			l286:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							{
								position289, tokenIndex289, depth289 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l290
								}
								position++
								goto l289
							l290:
								position, tokenIndex, depth = position289, tokenIndex289, depth289
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l287
								}
								position++
							}
						l289:
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l287
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l287
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l287
							}
							position++
							break
						}
					}

					goto l286
				l287:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
				}
				depth--
				add(ruleidentifier, position284)
			}
			return true
		l283:
			position, tokenIndex, depth = position283, tokenIndex283, depth283
			return false
		},
		/* 29 fields <- <(((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' ') | (&(';') ';'))* field isp* (fsep isp* (fsep isp*)* field)* ((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' ') | (&(';') ';'))* !.)> */
		func() bool {
			position291, tokenIndex291, depth291 := position, tokenIndex, depth
			{
				position292 := position
				depth++
			l293:
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l294
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l294
							}
							position++
							break
						case ' ':
							if buffer[position] != rune(' ') {
								goto l294
							}
							position++
							break
						default:
							if buffer[position] != rune(';') {
								goto l294
							}
							position++
							break
						}
					}

					goto l293
				l294:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
				}
				if !_rules[rulefield]() {
					goto l291
				}
			l296:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l297
					}
					goto l296
				l297:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
				}
			l298:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l299
					}
				l300:
					{
						position301, tokenIndex301, depth301 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l301
						}
						goto l300
					l301:
						position, tokenIndex, depth = position301, tokenIndex301, depth301
					}
				l302:
					{
						position303, tokenIndex303, depth303 := position, tokenIndex, depth
						if !_rules[rulefsep]() {
							goto l303
						}
					l304:
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l305
							}
							goto l304
						l305:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
						}
						goto l302
					l303:
						position, tokenIndex, depth = position303, tokenIndex303, depth303
					}
					if !_rules[rulefield]() {
						goto l299
					}
					goto l298
				l299:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
				}
			l306:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l307
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l307
							}
							position++
							break
						case ' ':
							if buffer[position] != rune(' ') {
								goto l307
							}
							position++
							break
						default:
							if buffer[position] != rune(';') {
								goto l307
							}
							position++
							break
						}
					}

					goto l306
				l307:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
				}
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if !matchDot() {
						goto l309
					}
					goto l291
				l309:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
				}
				depth--
				add(rulefields, position292)
			}
			return true
		l291:
			position, tokenIndex, depth = position291, tokenIndex291, depth291
			return false
		},
		/* 30 fsep <- <(';' / '\n')> */
		func() bool {
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{
				position311 := position
				depth++
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune(';') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('\n') {
						goto l310
					}
					position++
				}
			l312:
				depth--
				add(rulefsep, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 31 field <- <(name (isp* ',' isp* name)* isp+ type isp* ('=' isp* expr)? Action14)> */
		func() bool {
			position314, tokenIndex314, depth314 := position, tokenIndex, depth
			{
				position315 := position
				depth++
				if !_rules[rulename]() {
					goto l314
				}
			l316:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
				l318:
					{
						position319, tokenIndex319, depth319 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l319
						}
						goto l318
					l319:
						position, tokenIndex, depth = position319, tokenIndex319, depth319
					}
					if buffer[position] != rune(',') {
						goto l317
					}
					position++
				l320:
					{
						position321, tokenIndex321, depth321 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l321
						}
						goto l320
					l321:
						position, tokenIndex, depth = position321, tokenIndex321, depth321
					}
					if !_rules[rulename]() {
						goto l317
					}
					goto l316
				l317:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
				}
				if !_rules[ruleisp]() {
					goto l314
				}
			l322:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l323
					}
					goto l322
				l323:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
				}
				if !_rules[ruletype]() {
					goto l314
				}
			l324:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
				}
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('=') {
						goto l326
					}
					position++
				l328:
					{
						position329, tokenIndex329, depth329 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l329
						}
						goto l328
					l329:
						position, tokenIndex, depth = position329, tokenIndex329, depth329
					}
					if !_rules[ruleexpr]() {
						goto l326
					}
					goto l327
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
			l327:
				if !_rules[ruleAction14]() {
					goto l314
				}
				depth--
				add(rulefield, position315)
			}
			return true
		l314:
			position, tokenIndex, depth = position314, tokenIndex314, depth314
			return false
		},
		/* 32 name <- <(<((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> Action15)> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				{
					position332 := position
					depth++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l330
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l330
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l330
							}
							position++
							break
						}
					}

				l333:
					{
						position334, tokenIndex334, depth334 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l334
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l334
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l334
								}
								position++
								break
							}
						}

						goto l333
					l334:
						position, tokenIndex, depth = position334, tokenIndex334, depth334
					}
					depth--
					add(rulePegText, position332)
				}
				if !_rules[ruleAction15]() {
					goto l330
				}
				depth--
				add(rulename, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 33 type <- <(chan / func / qname / sname / ((&('*') pointer) | (&('[') array) | (&('M' | 'm') map)))> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if !_rules[rulechan]() {
						goto l340
					}
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if !_rules[rulefunc]() {
						goto l341
					}
					goto l339
				l341:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if !_rules[ruleqname]() {
						goto l342
					}
					goto l339
				l342:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if !_rules[rulesname]() {
						goto l343
					}
					goto l339
				l343:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					{
						switch buffer[position] {
						case '*':
							if !_rules[rulepointer]() {
								goto l337
							}
							break
						case '[':
							if !_rules[rulearray]() {
								goto l337
							}
							break
						default:
							if !_rules[rulemap]() {
								goto l337
							}
							break
						}
					}

				}
			l339:
				depth--
				add(ruletype, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 34 sname <- <(<((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> Action16)> */
		func() bool {
			position345, tokenIndex345, depth345 := position, tokenIndex, depth
			{
				position346 := position
				depth++
				{
					position347 := position
					depth++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l345
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l345
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l345
							}
							position++
							break
						}
					}

				l348:
					{
						position349, tokenIndex349, depth349 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l349
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l349
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l349
								}
								position++
								break
							}
						}

						goto l348
					l349:
						position, tokenIndex, depth = position349, tokenIndex349, depth349
					}
					depth--
					add(rulePegText, position347)
				}
				if !_rules[ruleAction16]() {
					goto l345
				}
				depth--
				add(rulesname, position346)
			}
			return true
		l345:
			position, tokenIndex, depth = position345, tokenIndex345, depth345
			return false
		},
		/* 35 qname <- <(<(((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+ '.' ((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)> Action17)> */
		func() bool {
			position352, tokenIndex352, depth352 := position, tokenIndex, depth
			{
				position353 := position
				depth++
				{
					position354 := position
					depth++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l352
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l352
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l352
							}
							position++
							break
						}
					}

				l355:
					{
						position356, tokenIndex356, depth356 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l356
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l356
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l356
								}
								position++
								break
							}
						}

						goto l355
					l356:
						position, tokenIndex, depth = position356, tokenIndex356, depth356
					}
					if buffer[position] != rune('.') {
						goto l352
					}
					position++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l352
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l352
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l352
							}
							position++
							break
						}
					}

				l359:
					{
						position360, tokenIndex360, depth360 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l360
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l360
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l360
								}
								position++
								break
							}
						}

						goto l359
					l360:
						position, tokenIndex, depth = position360, tokenIndex360, depth360
					}
					depth--
					add(rulePegText, position354)
				}
				if !_rules[ruleAction17]() {
					goto l352
				}
				depth--
				add(ruleqname, position353)
			}
			return true
		l352:
			position, tokenIndex, depth = position352, tokenIndex352, depth352
			return false
		},
		/* 36 array <- <('[' ']' type Action18)> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				if buffer[position] != rune('[') {
					goto l363
				}
				position++
				if buffer[position] != rune(']') {
					goto l363
				}
				position++
				if !_rules[ruletype]() {
					goto l363
				}
				if !_rules[ruleAction18]() {
					goto l363
				}
				depth--
				add(rulearray, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 37 map <- <(('m' / 'M') ('a' / 'A') ('p' / 'P') '[' isp* keytype isp* ']' type Action19)> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('M') {
						goto l365
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('A') {
						goto l365
					}
					position++
				}
			l369:
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('P') {
						goto l365
					}
					position++
				}
			l371:
				if buffer[position] != rune('[') {
					goto l365
				}
				position++
			l373:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l374
					}
					goto l373
				l374:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
				}
				if !_rules[rulekeytype]() {
					goto l365
				}
			l375:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l376
					}
					goto l375
				l376:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
				}
				if buffer[position] != rune(']') {
					goto l365
				}
				position++
				if !_rules[ruletype]() {
					goto l365
				}
				if !_rules[ruleAction19]() {
					goto l365
				}
				depth--
				add(rulemap, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 38 chan <- <(('c' / 'C') ('h' / 'H') ('a' / 'A') ('n' / 'N') isp+ type Action20)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('C') {
						goto l377
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('H') {
						goto l377
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('A') {
						goto l377
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('N') {
						goto l377
					}
					position++
				}
			l385:
				if !_rules[ruleisp]() {
					goto l377
				}
			l387:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l388
					}
					goto l387
				l388:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
				}
				if !_rules[ruletype]() {
					goto l377
				}
				if !_rules[ruleAction20]() {
					goto l377
				}
				depth--
				add(rulechan, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 39 func <- <(('f' / 'F') ('u' / 'U') ('n' / 'N') ('c' / 'C') isp* '(' isp* (param isp* (',' isp* param)*)? ')' isp* type? Action21)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('F') {
						goto l389
					}
					position++
				}
			l391:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('U') {
						goto l389
					}
					position++
				}
			l393:
				{
					position395, tokenIndex395, depth395 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l396
					}
					position++
					goto l395
				l396:
					position, tokenIndex, depth = position395, tokenIndex395, depth395
					if buffer[position] != rune('N') {
						goto l389
					}
					position++
				}
			l395:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('C') {
						goto l389
					}
					position++
				}
			l397:
			l399:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l400
					}
					goto l399
				l400:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
				}
				if buffer[position] != rune('(') {
					goto l389
				}
				position++
			l401:
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l402
					}
					goto l401
				l402:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
				}
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if !_rules[ruleparam]() {
						goto l403
					}
				l405:
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l406
						}
						goto l405
					l406:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
					}
				l407:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l408
						}
						position++
					l409:
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l410
							}
							goto l409
						l410:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
						}
						if !_rules[ruleparam]() {
							goto l408
						}
						goto l407
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
					goto l404
				l403:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
				}
			l404:
				if buffer[position] != rune(')') {
					goto l389
				}
				position++
			l411:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l412
					}
					goto l411
				l412:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
				}
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if !_rules[ruletype]() {
						goto l413
					}
					goto l414
				l413:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
				}
			l414:
				if !_rules[ruleAction21]() {
					goto l389
				}
				depth--
				add(rulefunc, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 40 keytype <- <(type Action22)> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				if !_rules[ruletype]() {
					goto l415
				}
				if !_rules[ruleAction22]() {
					goto l415
				}
				depth--
				add(rulekeytype, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 41 pointer <- <('*' type Action23)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				if buffer[position] != rune('*') {
					goto l417
				}
				position++
				if !_rules[ruletype]() {
					goto l417
				}
				if !_rules[ruleAction23]() {
					goto l417
				}
				depth--
				add(rulepointer, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 42 captures <- <(isp* capture isp* (',' isp* capture isp*)* !.)> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
			l421:
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l422
					}
					goto l421
				l422:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
				}
				if !_rules[rulecapture]() {
					goto l419
				}
			l423:
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l424
					}
					goto l423
				l424:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
				}
			l425:
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l426
					}
					position++
				l427:
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l428
						}
						goto l427
					l428:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
					}
					if !_rules[rulecapture]() {
						goto l426
					}
				l429:
					{
						position430, tokenIndex430, depth430 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l430
						}
						goto l429
					l430:
						position, tokenIndex, depth = position430, tokenIndex430, depth430
					}
					goto l425
				l426:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
				}
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if !matchDot() {
						goto l431
					}
					goto l419
				l431:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
				}
				depth--
				add(rulecaptures, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 43 capture <- <(eventid isp* ':' handlername isp* mappings isp* tags Action24)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				if !_rules[ruleeventid]() {
					goto l432
				}
			l434:
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l435
					}
					goto l434
				l435:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
				}
				if buffer[position] != rune(':') {
					goto l432
				}
				position++
				if !_rules[rulehandlername]() {
					goto l432
				}
			l436:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l437
					}
					goto l436
				l437:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
				}
				if !_rules[rulemappings]() {
					goto l432
				}
			l438:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l439
					}
					goto l438
				l439:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
				}
				if !_rules[ruletags]() {
					goto l432
				}
				if !_rules[ruleAction24]() {
					goto l432
				}
				depth--
				add(rulecapture, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 44 handlername <- <(<identifier> Action25)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				{
					position442 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l440
					}
					depth--
					add(rulePegText, position442)
				}
				if !_rules[ruleAction25]() {
					goto l440
				}
				depth--
				add(rulehandlername, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 45 eventid <- <(<[a-z]+> Action26)> */
		func() bool {
			position443, tokenIndex443, depth443 := position, tokenIndex, depth
			{
				position444 := position
				depth++
				{
					position445 := position
					depth++
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l443
					}
					position++
				l446:
					{
						position447, tokenIndex447, depth447 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l447
						}
						position++
						goto l446
					l447:
						position, tokenIndex, depth = position447, tokenIndex447, depth447
					}
					depth--
					add(rulePegText, position445)
				}
				if !_rules[ruleAction26]() {
					goto l443
				}
				depth--
				add(ruleeventid, position444)
			}
			return true
		l443:
			position, tokenIndex, depth = position443, tokenIndex443, depth443
			return false
		},
		/* 46 mappings <- <('(' (isp* mapping isp* (',' isp* mapping isp*)*)? ')')?> */
		func() bool {
			{
				position449 := position
				depth++
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l450
					}
					position++
					{
						position452, tokenIndex452, depth452 := position, tokenIndex, depth
					l454:
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l455
							}
							goto l454
						l455:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
						}
						if !_rules[rulemapping]() {
							goto l452
						}
					l456:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l457
							}
							goto l456
						l457:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
						}
					l458:
						{
							position459, tokenIndex459, depth459 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l459
							}
							position++
						l460:
							{
								position461, tokenIndex461, depth461 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l461
								}
								goto l460
							l461:
								position, tokenIndex, depth = position461, tokenIndex461, depth461
							}
							if !_rules[rulemapping]() {
								goto l459
							}
						l462:
							{
								position463, tokenIndex463, depth463 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l463
								}
								goto l462
							l463:
								position, tokenIndex, depth = position463, tokenIndex463, depth463
							}
							goto l458
						l459:
							position, tokenIndex, depth = position459, tokenIndex459, depth459
						}
						goto l453
					l452:
						position, tokenIndex, depth = position452, tokenIndex452, depth452
					}
				l453:
					if buffer[position] != rune(')') {
						goto l450
					}
					position++
					goto l451
				l450:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
				}
			l451:
				depth--
				add(rulemappings, position449)
			}
			return true
		},
		/* 47 mapping <- <(mappingname isp* '=' isp* bound Action27)> */
		func() bool {
			position464, tokenIndex464, depth464 := position, tokenIndex, depth
			{
				position465 := position
				depth++
				if !_rules[rulemappingname]() {
					goto l464
				}
			l466:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l467
					}
					goto l466
				l467:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
				}
				if buffer[position] != rune('=') {
					goto l464
				}
				position++
			l468:
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l469
					}
					goto l468
				l469:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
				}
				if !_rules[rulebound]() {
					goto l464
				}
				if !_rules[ruleAction27]() {
					goto l464
				}
				depth--
				add(rulemapping, position465)
			}
			return true
		l464:
			position, tokenIndex, depth = position464, tokenIndex464, depth464
			return false
		},
		/* 48 mappingname <- <(<identifier> Action28)> */
		func() bool {
			position470, tokenIndex470, depth470 := position, tokenIndex, depth
			{
				position471 := position
				depth++
				{
					position472 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l470
					}
					depth--
					add(rulePegText, position472)
				}
				if !_rules[ruleAction28]() {
					goto l470
				}
				depth--
				add(rulemappingname, position471)
			}
			return true
		l470:
			position, tokenIndex, depth = position470, tokenIndex470, depth470
			return false
		},
		/* 49 tags <- <('{' isp* tag isp* (',' isp* tag isp*)* '}')?> */
		func() bool {
			{
				position474 := position
				depth++
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('{') {
						goto l475
					}
					position++
				l477:
					{
						position478, tokenIndex478, depth478 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l478
						}
						goto l477
					l478:
						position, tokenIndex, depth = position478, tokenIndex478, depth478
					}
					if !_rules[ruletag]() {
						goto l475
					}
				l479:
					{
						position480, tokenIndex480, depth480 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l480
						}
						goto l479
					l480:
						position, tokenIndex, depth = position480, tokenIndex480, depth480
					}
				l481:
					{
						position482, tokenIndex482, depth482 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l482
						}
						position++
					l483:
						{
							position484, tokenIndex484, depth484 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l484
							}
							goto l483
						l484:
							position, tokenIndex, depth = position484, tokenIndex484, depth484
						}
						if !_rules[ruletag]() {
							goto l482
						}
					l485:
						{
							position486, tokenIndex486, depth486 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l486
							}
							goto l485
						l486:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
						}
						goto l481
					l482:
						position, tokenIndex, depth = position482, tokenIndex482, depth482
					}
					if buffer[position] != rune('}') {
						goto l475
					}
					position++
					goto l476
				l475:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
				}
			l476:
				depth--
				add(ruletags, position474)
			}
			return true
		},
		/* 50 tag <- <(tagname ('(' (isp* tagarg isp* (',' isp* tagarg isp*)*)? ')')? Action29)> */
		func() bool {
			position487, tokenIndex487, depth487 := position, tokenIndex, depth
			{
				position488 := position
				depth++
				if !_rules[ruletagname]() {
					goto l487
				}
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l489
					}
					position++
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
					l493:
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l494
							}
							goto l493
						l494:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
						}
						if !_rules[ruletagarg]() {
							goto l491
						}
					l495:
						{
							position496, tokenIndex496, depth496 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l496
							}
							goto l495
						l496:
							position, tokenIndex, depth = position496, tokenIndex496, depth496
						}
					l497:
						{
							position498, tokenIndex498, depth498 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l498
							}
							position++
						l499:
							{
								position500, tokenIndex500, depth500 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l500
								}
								goto l499
							l500:
								position, tokenIndex, depth = position500, tokenIndex500, depth500
							}
							if !_rules[ruletagarg]() {
								goto l498
							}
						l501:
							{
								position502, tokenIndex502, depth502 := position, tokenIndex, depth
								if !_rules[ruleisp]() {
									goto l502
								}
								goto l501
							l502:
								position, tokenIndex, depth = position502, tokenIndex502, depth502
							}
							goto l497
						l498:
							position, tokenIndex, depth = position498, tokenIndex498, depth498
						}
						goto l492
					l491:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
					}
				l492:
					if buffer[position] != rune(')') {
						goto l489
					}
					position++
					goto l490
				l489:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
				}
			l490:
				if !_rules[ruleAction29]() {
					goto l487
				}
				depth--
				add(ruletag, position488)
			}
			return true
		l487:
			position, tokenIndex, depth = position487, tokenIndex487, depth487
			return false
		},
		/* 51 tagname <- <(<identifier> Action30)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				{
					position505 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l503
					}
					depth--
					add(rulePegText, position505)
				}
				if !_rules[ruleAction30]() {
					goto l503
				}
				depth--
				add(ruletagname, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 52 tagarg <- <(<identifier> Action31)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l506
					}
					depth--
					add(rulePegText, position508)
				}
				if !_rules[ruleAction31]() {
					goto l506
				}
				depth--
				add(ruletagarg, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 53 for <- <(isp* forVar isp* (',' isp* forVar isp*)? (':' '=') isp* (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) isp+ expr isp* !.)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
			l511:
				{
					position512, tokenIndex512, depth512 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l512
					}
					goto l511
				l512:
					position, tokenIndex, depth = position512, tokenIndex512, depth512
				}
				if !_rules[ruleforVar]() {
					goto l509
				}
			l513:
				{
					position514, tokenIndex514, depth514 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l514
					}
					goto l513
				l514:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
				}
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l515
					}
					position++
				l517:
					{
						position518, tokenIndex518, depth518 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l518
						}
						goto l517
					l518:
						position, tokenIndex, depth = position518, tokenIndex518, depth518
					}
					if !_rules[ruleforVar]() {
						goto l515
					}
				l519:
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l520
						}
						goto l519
					l520:
						position, tokenIndex, depth = position520, tokenIndex520, depth520
					}
					goto l516
				l515:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
				}
			l516:
				if buffer[position] != rune(':') {
					goto l509
				}
				position++
				if buffer[position] != rune('=') {
					goto l509
				}
				position++
			l521:
				{
					position522, tokenIndex522, depth522 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l522
					}
					goto l521
				l522:
					position, tokenIndex, depth = position522, tokenIndex522, depth522
				}
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('R') {
						goto l509
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('A') {
						goto l509
					}
					position++
				}
			l525:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('N') {
						goto l509
					}
					position++
				}
			l527:
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l530
					}
					position++
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if buffer[position] != rune('G') {
						goto l509
					}
					position++
				}
			l529:
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l532
					}
					position++
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if buffer[position] != rune('E') {
						goto l509
					}
					position++
				}
			l531:
				if !_rules[ruleisp]() {
					goto l509
				}
			l533:
				{
					position534, tokenIndex534, depth534 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l534
					}
					goto l533
				l534:
					position, tokenIndex, depth = position534, tokenIndex534, depth534
				}
				if !_rules[ruleexpr]() {
					goto l509
				}
			l535:
				{
					position536, tokenIndex536, depth536 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l536
					}
					goto l535
				l536:
					position, tokenIndex, depth = position536, tokenIndex536, depth536
				}
				{
					position537, tokenIndex537, depth537 := position, tokenIndex, depth
					if !matchDot() {
						goto l537
					}
					goto l509
				l537:
					position, tokenIndex, depth = position537, tokenIndex537, depth537
				}
				depth--
				add(rulefor, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 54 forVar <- <(<identifier> Action32)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l538
					}
					depth--
					add(rulePegText, position540)
				}
				if !_rules[ruleAction32]() {
					goto l538
				}
				depth--
				add(ruleforVar, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 55 handlers <- <(isp* (fsep isp*)* handler isp* ((fsep isp*)+ handler isp*)* (fsep isp*)* !.)> */
		func() bool {
			position541, tokenIndex541, depth541 := position, tokenIndex, depth
			{
				position542 := position
				depth++
			l543:
				{
					position544, tokenIndex544, depth544 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l544
					}
					goto l543
				l544:
					position, tokenIndex, depth = position544, tokenIndex544, depth544
				}
			l545:
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l546
					}
				l547:
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l548
						}
						goto l547
					l548:
						position, tokenIndex, depth = position548, tokenIndex548, depth548
					}
					goto l545
				l546:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
				}
				if !_rules[rulehandler]() {
					goto l541
				}
			l549:
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l550
					}
					goto l549
				l550:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
				}
			l551:
				{
					position552, tokenIndex552, depth552 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l552
					}
				l555:
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l556
						}
						goto l555
					l556:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
					}
				l553:
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						if !_rules[rulefsep]() {
							goto l554
						}
					l557:
						{
							position558, tokenIndex558, depth558 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l558
							}
							goto l557
						l558:
							position, tokenIndex, depth = position558, tokenIndex558, depth558
						}
						goto l553
					l554:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
					}
					if !_rules[rulehandler]() {
						goto l552
					}
				l559:
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l560
						}
						goto l559
					l560:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
					}
					goto l551
				l552:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
				}
			l561:
				{
					position562, tokenIndex562, depth562 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l562
					}
				l563:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l564
						}
						goto l563
					l564:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
					}
					goto l561
				l562:
					position, tokenIndex, depth = position562, tokenIndex562, depth562
				}
				{
					position565, tokenIndex565, depth565 := position, tokenIndex, depth
					if !matchDot() {
						goto l565
					}
					goto l541
				l565:
					position, tokenIndex, depth = position565, tokenIndex565, depth565
				}
				depth--
				add(rulehandlers, position542)
			}
			return true
		l541:
			position, tokenIndex, depth = position541, tokenIndex541, depth541
			return false
		},
		/* 56 handler <- <(handlername '(' isp* (param isp* (',' isp* param isp*)*)? ')' (isp* type)? Action33)> */
		func() bool {
			position566, tokenIndex566, depth566 := position, tokenIndex, depth
			{
				position567 := position
				depth++
				if !_rules[rulehandlername]() {
					goto l566
				}
				if buffer[position] != rune('(') {
					goto l566
				}
				position++
			l568:
				{
					position569, tokenIndex569, depth569 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l569
					}
					goto l568
				l569:
					position, tokenIndex, depth = position569, tokenIndex569, depth569
				}
				{
					position570, tokenIndex570, depth570 := position, tokenIndex, depth
					if !_rules[ruleparam]() {
						goto l570
					}
				l572:
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l573
						}
						goto l572
					l573:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
					}
				l574:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l575
						}
						position++
					l576:
						{
							position577, tokenIndex577, depth577 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l577
							}
							goto l576
						l577:
							position, tokenIndex, depth = position577, tokenIndex577, depth577
						}
						if !_rules[ruleparam]() {
							goto l575
						}
					l578:
						{
							position579, tokenIndex579, depth579 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l579
							}
							goto l578
						l579:
							position, tokenIndex, depth = position579, tokenIndex579, depth579
						}
						goto l574
					l575:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
					}
					goto l571
				l570:
					position, tokenIndex, depth = position570, tokenIndex570, depth570
				}
			l571:
				if buffer[position] != rune(')') {
					goto l566
				}
				position++
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
				l582:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l583
						}
						goto l582
					l583:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
					}
					if !_rules[ruletype]() {
						goto l580
					}
					goto l581
				l580:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
				}
			l581:
				if !_rules[ruleAction33]() {
					goto l566
				}
				depth--
				add(rulehandler, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 57 paramname <- <(<identifier> Action34)> */
		func() bool {
			position584, tokenIndex584, depth584 := position, tokenIndex, depth
			{
				position585 := position
				depth++
				{
					position586 := position
					depth++
					if !_rules[ruleidentifier]() {
						goto l584
					}
					depth--
					add(rulePegText, position586)
				}
				if !_rules[ruleAction34]() {
					goto l584
				}
				depth--
				add(ruleparamname, position585)
			}
			return true
		l584:
			position, tokenIndex, depth = position584, tokenIndex584, depth584
			return false
		},
		/* 58 param <- <(paramname isp+ type Action35)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				if !_rules[ruleparamname]() {
					goto l587
				}
				if !_rules[ruleisp]() {
					goto l587
				}
			l589:
				{
					position590, tokenIndex590, depth590 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l590
					}
					goto l589
				l590:
					position, tokenIndex, depth = position590, tokenIndex590, depth590
				}
				if !_rules[ruletype]() {
					goto l587
				}
				if !_rules[ruleAction35]() {
					goto l587
				}
				depth--
				add(ruleparam, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 59 cparams <- <(isp* (cparam isp* (',' isp* cparam isp*)*)? !.)> */
		func() bool {
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
			l593:
				{
					position594, tokenIndex594, depth594 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l594
					}
					goto l593
				l594:
					position, tokenIndex, depth = position594, tokenIndex594, depth594
				}
				{
					position595, tokenIndex595, depth595 := position, tokenIndex, depth
					if !_rules[rulecparam]() {
						goto l595
					}
				l597:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l598
						}
						goto l597
					l598:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
					}
				l599:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l600
						}
						position++
					l601:
						{
							position602, tokenIndex602, depth602 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l602
							}
							goto l601
						l602:
							position, tokenIndex, depth = position602, tokenIndex602, depth602
						}
						if !_rules[rulecparam]() {
							goto l600
						}
					l603:
						{
							position604, tokenIndex604, depth604 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l604
							}
							goto l603
						l604:
							position, tokenIndex, depth = position604, tokenIndex604, depth604
						}
						goto l599
					l600:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
					}
					goto l596
				l595:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
				}
			l596:
				{
					position605, tokenIndex605, depth605 := position, tokenIndex, depth
					if !matchDot() {
						goto l605
					}
					goto l591
				l605:
					position, tokenIndex, depth = position605, tokenIndex605, depth605
				}
				depth--
				add(rulecparams, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 60 cparam <- <((var isp+)? tagname isp+ type Action36)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				{
					position608, tokenIndex608, depth608 := position, tokenIndex, depth
					if !_rules[rulevar]() {
						goto l608
					}
					if !_rules[ruleisp]() {
						goto l608
					}
				l610:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l611
						}
						goto l610
					l611:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
					}
					goto l609
				l608:
					position, tokenIndex, depth = position608, tokenIndex608, depth608
				}
			l609:
				if !_rules[ruletagname]() {
					goto l606
				}
				if !_rules[ruleisp]() {
					goto l606
				}
			l612:
				{
					position613, tokenIndex613, depth613 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l613
					}
					goto l612
				l613:
					position, tokenIndex, depth = position613, tokenIndex613, depth613
				}
				if !_rules[ruletype]() {
					goto l606
				}
				if !_rules[ruleAction36]() {
					goto l606
				}
				depth--
				add(rulecparam, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 61 var <- <(('v' / 'V') ('a' / 'A') ('r' / 'R') Action37)> */
		func() bool {
			position614, tokenIndex614, depth614 := position, tokenIndex, depth
			{
				position615 := position
				depth++
				{
					position616, tokenIndex616, depth616 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l617
					}
					position++
					goto l616
				l617:
					position, tokenIndex, depth = position616, tokenIndex616, depth616
					if buffer[position] != rune('V') {
						goto l614
					}
					position++
				}
			l616:
				{
					position618, tokenIndex618, depth618 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l619
					}
					position++
					goto l618
				l619:
					position, tokenIndex, depth = position618, tokenIndex618, depth618
					if buffer[position] != rune('A') {
						goto l614
					}
					position++
				}
			l618:
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l621
					}
					position++
					goto l620
				l621:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if buffer[position] != rune('R') {
						goto l614
					}
					position++
				}
			l620:
				if !_rules[ruleAction37]() {
					goto l614
				}
				depth--
				add(rulevar, position615)
			}
			return true
		l614:
			position, tokenIndex, depth = position614, tokenIndex614, depth614
			return false
		},
		/* 62 args <- <(isp* arg isp* (',' isp* arg isp*)* !.)> */
		func() bool {
			position622, tokenIndex622, depth622 := position, tokenIndex, depth
			{
				position623 := position
				depth++
			l624:
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l625
					}
					goto l624
				l625:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
				}
				if !_rules[rulearg]() {
					goto l622
				}
			l626:
				{
					position627, tokenIndex627, depth627 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l627
					}
					goto l626
				l627:
					position, tokenIndex, depth = position627, tokenIndex627, depth627
				}
			l628:
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l629
					}
					position++
				l630:
					{
						position631, tokenIndex631, depth631 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l631
						}
						goto l630
					l631:
						position, tokenIndex, depth = position631, tokenIndex631, depth631
					}
					if !_rules[rulearg]() {
						goto l629
					}
				l632:
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l633
						}
						goto l632
					l633:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
					}
					goto l628
				l629:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
				}
				{
					position634, tokenIndex634, depth634 := position, tokenIndex, depth
					if !matchDot() {
						goto l634
					}
					goto l622
				l634:
					position, tokenIndex, depth = position634, tokenIndex634, depth634
				}
				depth--
				add(ruleargs, position623)
			}
			return true
		l622:
			position, tokenIndex, depth = position622, tokenIndex622, depth622
			return false
		},
		/* 63 arg <- <(expr Action38)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				if !_rules[ruleexpr]() {
					goto l635
				}
				if !_rules[ruleAction38]() {
					goto l635
				}
				depth--
				add(rulearg, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 64 imports <- <(isp* (fsep isp*)* import isp* (fsep isp* (fsep isp*)* import isp*)* (fsep isp*)* !.)> */
		func() bool {
			position637, tokenIndex637, depth637 := position, tokenIndex, depth
			{
				position638 := position
				depth++
			l639:
				{
					position640, tokenIndex640, depth640 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l640
					}
					goto l639
				l640:
					position, tokenIndex, depth = position640, tokenIndex640, depth640
				}
			l641:
				{
					position642, tokenIndex642, depth642 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l642
					}
				l643:
					{
						position644, tokenIndex644, depth644 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l644
						}
						goto l643
					l644:
						position, tokenIndex, depth = position644, tokenIndex644, depth644
					}
					goto l641
				l642:
					position, tokenIndex, depth = position642, tokenIndex642, depth642
				}
				if !_rules[ruleimport]() {
					goto l637
				}
			l645:
				{
					position646, tokenIndex646, depth646 := position, tokenIndex, depth
					if !_rules[ruleisp]() {
						goto l646
					}
					goto l645
				l646:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
				}
			l647:
				{
					position648, tokenIndex648, depth648 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l648
					}
				l649:
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l650
						}
						goto l649
					l650:
						position, tokenIndex, depth = position650, tokenIndex650, depth650
					}
				l651:
					{
						position652, tokenIndex652, depth652 := position, tokenIndex, depth
						if !_rules[rulefsep]() {
							goto l652
						}
					l653:
						{
							position654, tokenIndex654, depth654 := position, tokenIndex, depth
							if !_rules[ruleisp]() {
								goto l654
							}
							goto l653
						l654:
							position, tokenIndex, depth = position654, tokenIndex654, depth654
						}
						goto l651
					l652:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
					}
					if !_rules[ruleimport]() {
						goto l648
					}
				l655:
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l656
						}
						goto l655
					l656:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
					}
					goto l647
				l648:
					position, tokenIndex, depth = position648, tokenIndex648, depth648
				}
			l657:
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if !_rules[rulefsep]() {
						goto l658
					}
				l659:
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l660
						}
						goto l659
					l660:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
					}
					goto l657
				l658:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
				}
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if !matchDot() {
						goto l661
					}
					goto l637
				l661:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
				}
				depth--
				add(ruleimports, position638)
			}
			return true
		l637:
			position, tokenIndex, depth = position637, tokenIndex637, depth637
			return false
		},
		/* 65 import <- <((tagname isp+)? '"' <(!'"' .)*> '"' Action39)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664, tokenIndex664, depth664 := position, tokenIndex, depth
					if !_rules[ruletagname]() {
						goto l664
					}
					if !_rules[ruleisp]() {
						goto l664
					}
				l666:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if !_rules[ruleisp]() {
							goto l667
						}
						goto l666
					l667:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
					}
					goto l665
				l664:
					position, tokenIndex, depth = position664, tokenIndex664, depth664
				}
			l665:
				if buffer[position] != rune('"') {
					goto l662
				}
				position++
				{
					position668 := position
					depth++
				l669:
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						{
							position671, tokenIndex671, depth671 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l671
							}
							position++
							goto l670
						l671:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
						}
						if !matchDot() {
							goto l670
						}
						goto l669
					l670:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
					}
					depth--
					add(rulePegText, position668)
				}
				if buffer[position] != rune('"') {
					goto l662
				}
				position++
				if !_rules[ruleAction39]() {
					goto l662
				}
				depth--
				add(ruleimport, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 67 Action0 <- <{
			p.varMappings = append(p.varMappings,
				data.VariableMapping{Value: p.bv, Variable: p.goVal})
			p.goVal.Type = nil
			p.bv.IDs = nil
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 69 Action1 <- <{
			p.goVal.Name = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 70 Action2 <- <{
			p.goVal.Type = p.valuetype
			p.valuetype = nil
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 71 Action3 <- <{
			p.assignments = append(p.assignments, data.Assignment{Expression: p.expr,
				Target: p.bv})
			p.bv.IDs = nil
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 72 Action4 <- <{
			p.bv.Kind = data.BoundSelf
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 73 Action5 <- <{
			p.bv.Kind = data.BoundDataset
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 74 Action6 <- <{
			p.bv.Kind = data.BoundProperty
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 75 Action7 <- <{
			p.bv.Kind = data.BoundStyle
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 76 Action8 <- <{
			p.bv.Kind = data.BoundClass
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 77 Action9 <- <{
			p.bv.Kind = data.BoundFormValue
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 78 Action10 <- <{
			p.bv.Kind = data.BoundEventValue
			if len(p.bv.IDs) == 0 {
				p.bv.IDs = append(p.bv.IDs, "")
			}
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 79 Action11 <- <{
			p.bv.IDs = append(p.bv.IDs, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 80 Action12 <- <{
			p.bv.IDs = append(p.bv.IDs, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 81 Action13 <- <{
			p.expr = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 82 Action14 <- <{
			var expr *string
			if p.expr != "" {
				expr = new(string)
				*expr = p.expr
			}
			for _, name := range p.names {
				p.fields = append(p.fields, &data.Field{Name: name, Type: p.valuetype, DefaultValue: expr})
			}
			p.expr = ""
			p.valuetype = nil
			p.names = nil
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 83 Action15 <- <{
			p.names = append(p.names, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 84 Action16 <- <{
			switch name := buffer[begin:end]; name {
			case "int":
				p.valuetype = &data.ParamType{Kind: data.IntType}
			case "bool":
				p.valuetype = &data.ParamType{Kind: data.BoolType}
			case "string":
				p.valuetype = &data.ParamType{Kind: data.StringType}
			default:
				p.valuetype = &data.ParamType{Kind: data.NamedType, Name: name}
			}
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 85 Action17 <- <{
			name := buffer[begin:end]
			if name == "js.Value" {
				p.valuetype = &data.ParamType{Kind: data.JSValueType}
			} else {
				p.valuetype = &data.ParamType{Kind: data.NamedType, Name: name}
			}
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 86 Action18 <- <{
			p.valuetype = &data.ParamType{Kind: data.ArrayType, ValueType: p.valuetype}
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 87 Action19 <- <{
			p.valuetype = &data.ParamType{Kind: data.MapType, KeyType: p.keytype, ValueType: p.valuetype}
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 88 Action20 <- <{
			p.valuetype = &data.ParamType{Kind: data.ChanType, ValueType: p.valuetype}
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 89 Action21 <- <{
			p.valuetype = &data.ParamType{Kind: data.FuncType, ValueType: p.valuetype,
				Params: p.params}
			p.params = nil
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 90 Action22 <- <{
			p.keytype = p.valuetype
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 91 Action23 <- <{
			p.valuetype = &data.ParamType{Kind: data.PointerType, ValueType: p.valuetype}
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 92 Action24 <- <{
			p.eventMappings = append(p.eventMappings, data.UnboundEventMapping{
				Event: p.expr, Handler: p.handlername, ParamMappings: p.paramMappings,
				Handling: p.eventHandling})
			p.eventHandling = data.AutoPreventDefault
			p.expr = ""
			p.paramMappings = make(map[string]data.BoundValue)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 93 Action25 <- <{
			p.handlername = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 94 Action26 <- <{
			p.expr = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 95 Action27 <- <{
			if _, ok := p.paramMappings[p.tagname]; ok {
				p.err = errors.New("duplicate param: " + p.tagname)
				return
			}
			p.paramMappings[p.tagname] = p.bv
			p.bv.IDs = nil
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 96 Action28 <- <{
			p.tagname = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 97 Action29 <- <{
			switch p.tagname {
			case "preventDefault":
				if p.eventHandling != data.AutoPreventDefault {
					p.err = errors.New("duplicate preventDefault")
					return
				}
				switch len(p.names) {
				case 0:
					p.eventHandling = data.PreventDefault
				case 1:
					switch p.names[0] {
					case "true":
						p.eventHandling = data.PreventDefault
					case "false":
						p.eventHandling = data.DontPreventDefault
					case "ask":
						p.eventHandling = data.AskPreventDefault
					default:
						p.err = fmt.Errorf("unsupported value for preventDefault: %s", p.names[0])
						return
					}
				default:
					p.err = errors.New("too many parameters for preventDefault")
					return
				}
			default:
				p.err = errors.New("unknown tag: " + p.tagname)
				return
			}
			p.names = nil
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 98 Action30 <- <{
			p.tagname = buffer[begin:end]
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 99 Action31 <- <{
			p.names = append(p.names, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 100 Action32 <- <{
			p.names = append(p.names, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 101 Action33 <- <{
			p.handlers = append(p.handlers, HandlerSpec{
				Name: p.handlername, Params: p.params, Returns: p.valuetype})
			p.valuetype = nil
			p.params = nil
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 102 Action34 <- <{
			p.paramnames = append(p.paramnames, buffer[begin:end])
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 103 Action35 <- <{
			name := p.paramnames[len(p.paramnames)-1]
			p.paramnames = p.paramnames[:len(p.paramnames)-1]
			for _, para := range p.params {
				if para.Name == name {
					p.err = errors.New("duplicate param name: " + para.Name)
					return
				}
			}

			p.params = append(p.params, data.Param{Name: name, Type: p.valuetype})
			p.valuetype = nil
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 104 Action36 <- <{
			p.cParams = append(p.cParams, data.ComponentParam{
				Name: p.tagname, Type: *p.valuetype, IsVar: p.isVar})
			p.valuetype = nil
			p.isVar = false
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 105 Action37 <- <{
			p.isVar = true
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 106 Action38 <- <{
		  p.names = append(p.names, p.expr)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 107 Action39 <- <{
			path := buffer[begin:end]
			if p.tagname == "" {
				lastDot := strings.LastIndexByte(path, '/')
				if lastDot == -1 {
					p.tagname = path
				} else {
					p.tagname = path[lastDot+1:]
				}
			}
			if _, ok := p.imports[p.tagname]; ok {
				p.err = errors.New("duplicate import name: " + p.tagname)
				return
			}
			p.imports[p.tagname] = path
			p.tagname = ""
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
	}
	p.rules = _rules
}
