package unityai

import (
	"github.com/tofuhua/unityai/format"
)

type NavMeshData struct {
	m_NavMeshBuildSettings NavMeshBuildSettings
	m_NavMeshTiles         []NavMeshTileData
	m_HeightMeshes         []HeightMeshData
	m_OffMeshLinks         []AutoOffMeshLinkData
	m_FilterAreaCosts      []float32
	m_SourceBounds         AABB
	m_Rotation             Quaternionf
	m_Position             Vector3f
	m_AgentTypeID          int32
}

func (this *NavMeshData) GetPosition() Vector3f {
	return this.m_Position
}

func (this *NavMeshData) GetRotation() Quaternionf {
	return this.m_Rotation
}

func (this *NavMeshData) GetNavMeshTiles() []NavMeshTileData {
	return this.m_NavMeshTiles
}

func (this *NavMeshData) GetNavMeshBuildSettings() NavMeshBuildSettings {
	return this.m_NavMeshBuildSettings
}

func (this *NavMeshData) GetOffMeshLinks() []AutoOffMeshLinkData {
	return this.m_OffMeshLinks
}

func (this *NavMeshData) GetFilterAreaCosts() []float32 {
	return this.m_FilterAreaCosts
}

func (this *NavMeshData) AddOffMeshLink(source, target Vector3f, area, linkDirection byte) {
	this.m_OffMeshLinks = append(this.m_OffMeshLinks,
		AutoOffMeshLinkData{m_Start: source, m_End: target,
			m_Area: area, m_LinkDirection: linkDirection})
}

func (this *NavMeshData) GetAgentTypeId() int32 {
	return this.m_AgentTypeID
}

func (this *NavMeshData) Clone() *NavMeshData {
	data := NavMeshData{}
	data.m_NavMeshBuildSettings = this.m_NavMeshBuildSettings
	data.m_SourceBounds = this.m_SourceBounds
	data.m_Rotation = this.m_Rotation
	data.m_Position = this.m_Position
	data.m_AgentTypeID = this.m_AgentTypeID

	data.m_NavMeshTiles = make([]NavMeshTileData, len(this.m_NavMeshTiles))
	for i := range this.m_NavMeshTiles {
		otile := this.m_NavMeshTiles[i]
		ntile := &data.m_NavMeshTiles[i]
		ntile.m_Hash = otile.m_Hash
		ntile.m_MeshData = make([]byte, len(otile.m_MeshData))
		copy(ntile.m_MeshData, otile.m_MeshData)
	}

	data.m_HeightMeshes = make([]HeightMeshData, len(this.m_HeightMeshes))
	for i := range this.m_HeightMeshes {
		omesh := this.m_HeightMeshes[i]
		nmesh := &this.m_HeightMeshes[i]
		nmesh.m_Bounds = omesh.m_Bounds
		nmesh.m_Vertices = make([]Vector3f, len(omesh.m_Vertices))
		copy(nmesh.m_Vertices, omesh.m_Vertices)
		nmesh.m_Indices = make([]int32, len(omesh.m_Indices))
		copy(nmesh.m_Indices, omesh.m_Indices)
		nmesh.m_Nodes = make([]HeightMeshBVNode, len(omesh.m_Nodes))
		copy(nmesh.m_Nodes, omesh.m_Nodes)
	}

	data.m_OffMeshLinks = make([]AutoOffMeshLinkData, len(this.m_OffMeshLinks))
	copy(data.m_OffMeshLinks, this.m_OffMeshLinks)

	data.m_FilterAreaCosts = make([]float32, len(this.m_FilterAreaCosts))
	copy(data.m_FilterAreaCosts, this.m_FilterAreaCosts)
	return &data
}

type NavMeshTileData struct {
	m_MeshData []byte
	m_Hash     [16]byte
}

func (this NavMeshTileData) GetData() []byte {
	return this.m_MeshData
}

type HeightMeshData struct {
	m_Vertices []Vector3f
	m_Indices  []int32
	m_Nodes    []HeightMeshBVNode
	m_Bounds   AABB
}

type AutoOffMeshLinkData struct {
	m_Start         Vector3f
	m_End           Vector3f
	m_Radius        float32
	m_LinkType      uint16 // Off-mesh poly flags.
	m_Area          byte   // Off-mesh poly  area ids.
	m_LinkDirection byte   // Off-mesh connection direction flags (NavMeshLinkDirectionFlags)
}

type HeightMeshBVNode struct {
	min, max Vector3f
	i, n     int32
}

func fromVector3f(data format.Vector3f) Vector3f {
	return Vector3f{data.X, data.Y, data.Z}
}

func fromAABB(data format.AABB) AABB {
	return AABB{
		fromVector3f(data.M_Center), fromVector3f(data.M_Extent),
	}
}

func fromHeightMeshBVNode(data *format.HeightMeshBVNode) HeightMeshBVNode {
	return HeightMeshBVNode{
		fromVector3f(data.Min), fromVector3f(data.Max), data.I, data.N,
	}
}

func fromHeightMeshData(data *format.HeightMeshData) HeightMeshData {
	var nvData HeightMeshData
	nvData.m_Indices = data.M_Indices
	nvData.m_Bounds = fromAABB(data.M_Bounds)
	nvData.m_Nodes = make([]HeightMeshBVNode, len(data.M_Nodes), len(data.M_Nodes))
	for i := range data.M_Nodes {
		nvData.m_Nodes[i] = fromHeightMeshBVNode(&data.M_Nodes[i])
	}
	nvData.m_Vertices = make([]Vector3f, len(data.M_Vertices), len(data.M_Vertices))
	for i := range data.M_Vertices {
		nvData.m_Vertices[i] = fromVector3f(data.M_Vertices[i])
	}
	return nvData
}

func fromNavMeshBuildSetting(data *format.NavMeshBuildSettings) NavMeshBuildSettings {
	return NavMeshBuildSettings{
		data.AgentTypeID, data.AgentRadius, data.AgentHeight,
		data.AgentSlope, data.AgentClimb, data.LedgeDropHeight,
		data.MaxJumpAcrossDistance, data.MinRegionArea, data.ManualCellSize,
		data.CellSize, data.ManualTileSize, data.TileSize, data.AccuratePlacement,
	}
}

func fromQuaterionf(data format.Quaternionf) Quaternionf {
	return Quaternionf{
		data.X, data.Y, data.Z, data.W,
	}
}
func fromNavMeshTileData(data *format.NavMeshTileData) NavMeshTileData {
	return NavMeshTileData{
		data.M_MeshData, data.M_Hash,
	}
}

func fromOffMeshLinkData(data *format.AutoOffMeshLinkData) AutoOffMeshLinkData {
	var linkDirection byte
	if data.M_LinkDirection {
		linkDirection = 1
	}
	return AutoOffMeshLinkData{
		fromVector3f(data.M_Start), fromVector3f(data.M_End), data.M_Radius,
		data.M_LinkType, data.M_Area, linkDirection,
	}
}

func NewDataFromFormat(data *format.NavMeshData) *NavMeshData {
	var nvData NavMeshData
	nvData.m_AgentTypeID = data.M_AgentTypeID
	nvData.m_Position = fromVector3f(data.M_Position)
	nvData.m_NavMeshBuildSettings = fromNavMeshBuildSetting(&data.M_NavMeshBuildSettings)
	nvData.m_SourceBounds = fromAABB(data.M_SourceBounds)
	nvData.m_Rotation = fromQuaterionf(data.M_Rotation)
	nvData.m_HeightMeshes = make([]HeightMeshData, len(data.M_HeightMeshes), len(data.M_HeightMeshes))
	for i := range data.M_HeightMeshes {
		nvData.m_HeightMeshes[i] = fromHeightMeshData(&data.M_HeightMeshes[i])
	}
	nvData.m_NavMeshTiles = make([]NavMeshTileData, len(data.M_NavMeshTiles), len(data.M_NavMeshTiles))
	for i := range data.M_NavMeshTiles {
		nvData.m_NavMeshTiles[i] = fromNavMeshTileData(&data.M_NavMeshTiles[i])
	}
	nvData.m_OffMeshLinks = make([]AutoOffMeshLinkData, len(data.M_AdditionalData.OffMeshLinks), len(data.M_AdditionalData.OffMeshLinks))
	for i := range data.M_AdditionalData.OffMeshLinks {
		nvData.m_OffMeshLinks[i] = fromOffMeshLinkData(&data.M_AdditionalData.OffMeshLinks[i])
	}
	nvData.m_FilterAreaCosts = make([]float32, kMaxAreas, kMaxAreas)
	copy(nvData.m_FilterAreaCosts, data.M_AdditionalData.AreaCosts)
	return &nvData
}

func NewNavMeshObstacleFromFormat(obsInfo format.SceneObsData) *NavMeshObstacle {
	position := fromVector3f(obsInfo.Position)
	scale := fromVector3f(obsInfo.Scale)
	rotation := fromQuaterionf(obsInfo.Rotation)
	obs := NewNavMeshObstacle(NavMeshObstacleShape(obsInfo.Shape), position, scale, rotation)
	obs.SetCenter(fromVector3f(obsInfo.Center))
	if obsInfo.Shape == 1 {
		obs.SetSize(fromVector3f(obsInfo.Size))
	} else {
		obs.SetHeight(obsInfo.Height)
		obs.SetRadius(obsInfo.Radius)
	}
	return obs
}
