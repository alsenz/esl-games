package lesson2

import(
	"errors"
	"golang.org/x/net/html"
	"strings"
)

const (
	EachPlayerTagName string = "each-player"
	EachRoundTagName string = "each-round"
	EachTeamTagName string = "each-team"
	SelectPlayerTagName string = "select-player"
	SelectRoundTagName string = "select-round"
	SelectTeamTagName string = "select-team"
)
var loopingTags = map[string]bool{
	EachPlayerTagName: true,
	EachRoundTagName: true,
	EachTeamTagName: true,
}
var filteringTags = map[string]bool {
	SelectPlayerTagName: true,
	SelectRoundTagName: true,
	SelectTeamTagName: true
}

type Renderer struct {
	model *Model
	slideTemplate string
}

func NewRenderer(mdl *Model, slideTemplate string) *Renderer {
	return &Renderer{mdl, slideTemplate}
}

func deepCopyNodeReTag(n *html.Node, reTag string) *html.Node {
	ret := deepCopyNode(n)
	ret.Data = reTag
}

// Returns an unparented copy of this node. If reTag is provided, the root element tag will be converted to this tag
func deepCopyNode(n *html.Node) *html.Node {
	cpy := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}
	copy(cpy.Attr, n.Attr)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		cpy.AppendChild(deepCopyNode(child))
	}
	return cpy
}

func clearChildren(n *html.Node) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		n.RemoveChild(child)
	}
}

func (renderer *Renderer) renderRecursive(mdl *Model, node *html.Node) error {
	if node.Type == html.TextNode {
		//TODO from here...
		//TODO pick this up - and check the return type
		node.Data = mdl.Eval(node.Data) //TODO let Eval sit on the model
		return nil
	} else if node.Type == html.ElementNode {
		if _, containsLooping := loopingTags[node.Data]; containsLooping {
			// Creates a copy of this node as a div
			protoDiv := deepCopyNodeReTag(node, "div")
			// Clears the contents of this node so we can repopulate
			clearChildren(node)
			//TODO step 3: for each group (of the right kind)
			//TODO TODO we now need to implement the switch and filter
			switch node.Data {
			case EachPlayerTagName:
				//TODO ... step 3 (contd) copy myself
				//TODO - implement on model
			case EachTeamTagName:
				//TODO ... step 3 (contd) copy myself
				//TODO - implement on model
			case EachRoundTagName:
				//TODO ... step 3 (contd) copy myself
				//TODO - implement on model
			default:
				return errors.New("Logic error: unimplemented looping tag: " + node.Data)
			}
			//TODO ... step 3 (contd) copy myself
			//TODO ... step 4 (contd) apply render recursive with model group
			//TODO ... step 5 (contd) add myself back to the parent
		} else if _, containsFiltering := loopingTags[node.Data]; containsFiltering {

		} else {
			//TODO
		}
	} else {
		return nil
	}
}

func (renderer *Renderer) Render() (string, error) {
	doc, err := html.Parse(strings.NewReader(renderer.slideTemplate))
	if err != nil {
		return "", err
	}
	// This recursion modifies the DOM in-place
	if err = renderer.renderRecursive(renderer.model, doc); err != nil {
		return "", err
	}
	//TODO make this back into a string... TODO TODO
	return "", nil
}


