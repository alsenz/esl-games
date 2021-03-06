package lesson2

import(
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"strconv"
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
	SelectTeamTagName: true,
}

const (
	ByAccumulatedRankingAttribute string = "by-accumulated-ranking"
	ByRankingAttribute string = "by-ranking"
)

type Renderer struct {
	model *ClientModel
	slideTemplate string
}

func NewRenderer(mdl *ClientModel, slideTemplate string) *Renderer {
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

// Returns the score to rank by and whether it's accumulated ranking
func rankingInfo(node *html.Node) (ScoreName, bool) {
	if node.Data != EachPlayerTagName {
		return "", false
	}
	for _, attr := range node.Attr {
		if attr.Key == ByRankingAttribute {
			return ScoreName(attr.Val), false
		} else if attr.Key == ByAccumulatedRankingAttribute {
			return ScoreName(attr.Val), true
		}
	}
	return "", false
}

func nameAttr(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "name" {
			return attr.Val
		}
	}
	return ""
}

func roundAttrs(node *html.Node) (Round, error) {
	round := Round{0, 0, 0}
	var err error
	for _, attr := range node.Attr {
		if strings.ToLower(attr.Key) == "act" {
			if round.Act, err = strconv.ParseUint(attr.Val, 10, 64); errr != nil {
				return round, err
			}
		} else if strings.ToLower(attr.Key) == "scene" {
			if round.Scene, err = strconv.ParseUint(attr.Val, 10, 64); errr != nil {
				return round, err
			}
		} else if strings.ToLower(attr.Key) == "rep" {
			if round.Rep, err = strconv.ParseUint(attr.Val, 10, 64); errr != nil {
				return round, err
			}
		}
	}
	return round, nil
}

func (renderer *Renderer) renderRecursive(mdl *ClientModel, node *html.Node) error {
	if node.Type == html.TextNode {
		var err error
		if node.Data, err = mdl.Eval(node.Data); err != nil {
			return err
		}
		return nil
	} else if node.Type == html.ElementNode {
		if _, containsLooping := loopingTags[node.Data]; containsLooping {
			// Creates a copy of this node as a div
			protoDiv := deepCopyNodeReTag(node, "div")
			// Clears the contents of this node so we can repopulate
			clearChildren(node)
			// Switch for each possible tag
			switch strings.ToLower(node.Data) {
			case EachPlayerTagName: // <each-player by-ranking="default"><h1>{{ .GetPlayers[0].Name }}</h1></each-player>
				//Check to see which version of ForEachPlayer we want
				if rankingScore, isAccumulated := rankingInfo(node);
					len(rankingScore) > 0 && isAccumulated {
					for _, subMdl := range mdl.ForEachPlayerByRank(rankingScore) {
						cpyDiv := deepCopyNode(protoDiv)
						renderer.renderRecursive(&subMdl, cpyDiv)
						node.AppendChild(cpyDiv)
					}
				} else if len(rankingScore) > 0 && !isAccumulated {
					for _, subMdl := range mdl.ForEachPlayerByAccumulatedRank(rankingScore) {
						cpyDiv := deepCopyNode(protoDiv)
						renderer.renderRecursive(&subMdl, cpyDiv)
						node.AppendChild(cpyDiv)
					}
				} else { //Bog standard player looping
					for _, subMdl := range mdl.ForEachPlayer() {
						cpyDiv := deepCopyNode(protoDiv)
						renderer.renderRecursive(&subMdl, cpyDiv)
						node.AppendChild(cpyDiv)
					}
				}
			case EachTeamTagName: //Fixme: may need team rankings on these
				for _, subMdl := range mdl.ForEachTeam() {
					cpyDiv := deepCopyNode(protoDiv)
					renderer.renderRecursive(&subMdl, cpyDiv)
					node.AppendChild(cpyDiv)
				}
			case EachRoundTagName:
				for _, subMdl := range mdl.ForEachRound() {
					cpyDiv := deepCopyNode(protoDiv)
					renderer.renderRecursive(&subMdl, cpyDiv)
					node.AppendChild(cpyDiv)
				}
			default:
				return errors.New("Logic error: unimplemented looping tag: " + node.Data)
			}
		} else if _, containsFiltering := loopingTags[node.Data]; containsFiltering {
			// Creates a copy of this node as a div
			cpyDiv := deepCopyNodeReTag(node, "div")
			// Clears the contents of this node so we can repopulate
			clearChildren(node)
			nodeName := nameAttr(node)
			var subMdl ClientModel
			switch strings.ToLower(node.Data) {
			case SelectPlayerTagName:
				subMdl = mdl.GetEntriesByPlayerName(nodeName)
			case SelectTeamTagName:
				subMdl = mdl.GetEntriesByTeam(nodeName)
			case SelectRoundTagName:
				round, err := roundAttrs(node)
				if err != nil {
					return err
				}
				subMdl = mdl.GetEntriesByRound(round.Act, round.Scene, round.Rep)
			}
			renderer.renderRecursive(&subMdl, cpyDiv)
			node.AppendChild(cpyDiv)
		} else { //Element node child recurse
			// Just loop over the children and modify in-place
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				renderer.renderRecursive(mdl, node)
			}
		}
	} else { //Other node child recurse
		// Just loop over the children and modify in-place
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			renderer.renderRecursive(mdl, node)
		}
	}
	return nil
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
	buf := &bytes.Buffer{}
	if err = html.Render(buf, doc); err != nil {
		return "", err
	}
	return buf.String(), nil
}


