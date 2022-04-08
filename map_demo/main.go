package main

import (
	"container/list"
	"fmt"
	util "github.com/tti6o/go-demo/pkg"
	"sync"
	"time"
)

//全局容器
var TagMetaBktGlobal map[string]*VirtualCTNR

//容器
type VirtualCTNR struct {
	tagMetaVerCTNR map[string]*TagMetaVersionBkt // {treeId: TagMetaVersionBkt}
}

//一棵树的单元
type TagMetaVersionBkt struct {
	TagMetaList    *list.List
	CurrentVersion string // 当前最新版本
	ActiveVersion  string // 处理中的版本，只会存在一个
	ITagMeta       ITagMeta
	//TagContStore   *TagContentStore         //收到查询请求时读取
	//ChangeContent  *TagContentStore         //接收执行事务队列时写入
	//CombTagCache   map[string][]interface{} // 组合标签存储表
	//CombTagMu      sync.RWMutex
	//SubTreeMap     map[string][]string // 某个节点下面所有子节点的深度优先遍历结果
	//SubTreeMu      sync.RWMutex
}

type TagMetaListNode struct {
	Version         string
	TagMetaSnapshot ITagMeta
	TagContRel      map[string][]interface{} // 存放修改之前的标签-内容表
	ContTagRel      map[string][]interface{} // 存储修改之前的内容-标签表
	ContStore       map[string]interface{}   // 存储修改之前的内容，当解绑时，需要将内容记录下来
}

type ITagMeta interface {
}

//树的元数据
type TagMeta struct {
	tagTree *PaasTree
}

//应用树
type PaasTree struct {
	*MultiwayTree
	pathMu      sync.RWMutex
	idPathMap   map[string]string
	namePathMap map[string]string
}

//多叉树
type MultiwayTree struct {
	ID           string
	Root         string
	TagNodeMap   map[string]*Node    // nodeID: *Node
	NodeChildren map[string][]string // nodeID: childrenID
	rwlock       sync.RWMutex
}

//树节点
type Node struct {
	ID         string
	Tag        string
	ParentID   []string //parents list
	Attr       map[string]interface{}
	CreateTime time.Time
	Creator    string
}

//type TagContentStore struct {
//	BaseTagCont   map[string][]interface{}
//	BaseContTag   map[string][]interface{}
//	BaseContStore map[string]interface{}
//	Lock          *sync.RWMutex
//}

func InitTagMetaBktGlobal() {
	TagMetaBktGlobal = make(map[string]*VirtualCTNR)
}

func NewContainer(treeId string) (container *VirtualCTNR, err error) {
	var newTagMetaRes ITagMeta
	if newTagMetaRes, err = NewTagMeta(treeId); err != nil {
		return nil, err
	}
	versionCTNR := TagMetaVersionBkt{
		ActiveVersion: "",
		TagMetaList:   list.New(),
		ITagMeta:      newTagMetaRes,
		//TagContStore:  model.NewTagContentStore(),
		//ChangeContent: model.NewTagContentStore(),
		//CombTagCache:  map[string][]interface{}{},
		//SubTreeMap:    map[string][]string{},
	}
	bktMap := make(map[string]*TagMetaVersionBkt, 1)
	bktMap[treeId] = &versionCTNR
	container = &VirtualCTNR{
		tagMetaVerCTNR: bktMap,
	}
	return
}

func (vb *VirtualCTNR) AddContainer(treeId, version string) (err error) {
	newBkt := vb.tagMetaVerCTNR[treeId]
	//if newBkt == nil {
	//	return errors.New("not found")
	//}
	//vb.SetTagContStore(newBkt)
	//newBkt.ChangeContent.BaseTagCont = nil
	//newBkt.ChangeContent.BaseContStore = nil
	//newBkt.ChangeContent.BaseContTag = nil
	//newBkt.ChangeContent = nil
	// 版本生成结束后，将该版本号存入该字段
	newBkt.CurrentVersion = version
	oldContainer := TagMetaBktGlobal[treeId]
	if oldContainer != nil {
		//TODO
		//oldBkt := oldContainer.tagMetaVerCTNR[treeId]
		//if oldBkt != nil {
		//	tagMetaList := newBkt.TagMetaList
		//	elem := tagMetaList.Back()
		//	if elem == nil {
		//		return
		//	}
		//	for oldNode := oldBkt.TagMetaList.Back(); oldNode != nil && tagMetaList.Len() <= 2; oldNode = oldNode.Prev() {
		//		tagMetaList.InsertBefore(oldNode, elem)
		//	}
		//}
		oldContainer.tagMetaVerCTNR = nil
		oldContainer = nil //todo
	}

	TagMetaBktGlobal[treeId] = vb
	return
}

func NewTagMeta(treeId string) (tg ITagMeta, err error) {
	var (
		paasTree *PaasTree
	)

	if paasTree, err = NewPaasTree(treeId); err != nil {
		paasTree = nil // clear
		return nil, err
	}
	tg = &TagMeta{
		tagTree: paasTree,
	}
	return
}

func NewPaasTree(treeId string) (pt *PaasTree, err error) {
	// only create root paas tree
	paasTree := NewMultiwayTree(treeId)
	// create root node, one paas tree has one root node
	if paasTree.Root == "" {
		if _, err = paasTree.CreateRootNode(treeId); err != nil {
			return nil, err
		}
	}

	pt = &PaasTree{
		MultiwayTree: paasTree,
		idPathMap:    make(map[string]string),
		namePathMap:  make(map[string]string),
	}

	return
}

func NewMultiwayTree(tid string) *MultiwayTree {
	// create a new tree
	tagTree := &MultiwayTree{
		ID:           tid,
		Root:         "",
		TagNodeMap:   make(map[string]*Node),
		NodeChildren: make(map[string][]string),
	}

	return tagTree
}

func (t *MultiwayTree) CreateRootNode(root string) (node *Node, err error) {
	t.Root = root
	node = NewNode(root, "root", "system")
	//node.SetAttr(nil)
	if err = t.AddNode(node, []string{}); err != nil {
		return nil, err
	}

	return
}

func NewNode(nid, tag, creator string) *Node {
	return &Node{
		ID:         nid,
		Tag:        tag,
		ParentID:   []string{},
		Attr:       make(map[string]interface{}),
		Creator:    creator,
		CreateTime: time.Now(),
	}
}

func (t *MultiwayTree) AddNode(node *Node, pidList []string) (err error) {
	t.TagNodeMap[node.ID] = node
	// update self parents
	//TODO
	//err = t.updateNodeParents(node.ID, pidList, ADD)
	//if err != nil {
	//	return err
	//}
	//// update node`s parents NodeChildren
	//err = t.updateNodeChildren(node.ID, pidList, ADD)
	//if err != nil {
	//	return err
	//}
	return
}

func main() {
	treeId := "tof4_org"
	util.PrintMemStats("InitTagMetaBktGlobal before")
	InitTagMetaBktGlobal()
	util.PrintMemStats("InitTagMetaBktGlobal after")
	for ver := 1; ver <= 10000000; ver++ {
		version := fmt.Sprintf("%d", ver)
		container, _ := NewContainer(treeId)
		//runtime.GC()
		container.AddContainer(treeId, version)
		//runtime.GC()
	}
	util.PrintMemStats("end")
}

//var NodeMap = make(map[int]*Node)
//
//var IntMap = make(map[int]int)
//
////var m map[int]interface{}
//
//func NewNode(num int) *Node {
//	m := make(map[int]int, 0)
//	m[num] = num
//	node := &Node{
//		ID:   num,
//		Attr: m,
//	}
//	return node
//}
//
//func TestIntMap() {
//	log.Printf("initMap before")
//	util.PrintMemStats("TestMapRelease-1")
//
//	var cnt = 100000
//	for i := 0; i < cnt; i++ {
//		IntMap[i] = i
//	}
//	log.Println(len(IntMap))
//
//	log.Printf("IntMap after")
//
//	util.PrintMemStats("TestMapRelease-2")
//	for i := 0; i < cnt-1; i++ {
//		delete(IntMap, i)
//	}
//
//	log.Println(len(IntMap))
//
//	log.Printf("delete after")
//
//	util.PrintMemStats("TestMapRelease-3")
//
//	IntMap = nil
//
//	log.Printf("set nil after")
//
//	util.PrintMemStats("TestMapRelease-4")
//}
//
//func TestNodeMap() {
//	log.Printf("NodeMap before")
//	util.PrintMemStats("TestMapRelease-1")
//
//	var cnt = 100000
//	for i := 0; i < cnt; i++ {
//		NodeMap[i] = NewNode(i)
//	}
//	log.Println(len(NodeMap))
//
//	log.Printf("NodeMap after")
//
//	util.PrintMemStats("TestMapRelease-2")
//	for i := 0; i < cnt-1; i++ {
//		delete(NodeMap, i)
//	}
//
//	log.Println(len(NodeMap))
//
//	log.Printf("delete after")
//
//	util.PrintMemStats("TestMapRelease-3")
//
//	NodeMap = nil
//
//	log.Printf("set nil after")
//
//	util.PrintMemStats("TestMapRelease-4")
//}
//
//func main() {
//	TestNodeMap()
//}
