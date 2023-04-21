package aibase

import "github.com/i4de/rulex/typex"

func NewBuildInAiRuntime(re typex.RuleX) typex.XAiBase {
	ai := new(BuildInAiRuntime)
	ai.re = re
	ai.aiBases = make(map[string]*typex.AiBase)
	ai.LoadAi(&typex.AiBase{
		UUID:     "MNIST",
		Name:     "MNIST example",
		Type:     "BUILDIN", // 内置类型
		Filepath: "_",
		Config: map[string]interface{}{
			"Layer": 10,
			"Alg":   "ANN",
		},
		AiBase:      ai,
		Description: "Classic ANN MNIST example",
	})
	return ai
}

/*
*
* 管理器
*
 */

type BuildInAiRuntime struct {
	re      typex.RuleX
	aiBases map[string]*typex.AiBase
}

func (ai *BuildInAiRuntime) GetRuleX() typex.RuleX {

	return ai.re
}
func (ai *BuildInAiRuntime) ListAi() []*typex.AiBase {
	result := []*typex.AiBase{}
	for _, v := range ai.aiBases {
		result = append(result, v)
	}
	return result

}
func (ai *BuildInAiRuntime) LoadAi(Ai *typex.AiBase) error {
	ai.aiBases[Ai.Name] = Ai
	return nil
}
func (ai *BuildInAiRuntime) GetAi(uuid string) *typex.AiBase {
	return ai.aiBases[uuid]
}
func (ai *BuildInAiRuntime) RemoveAi(uuid string) error {
	return nil
}
func (ai *BuildInAiRuntime) UpdateAi(Ai typex.AiBase) error {
	return nil
}
func (ai *BuildInAiRuntime) StartAi(uuid string) error {
	return nil
}
func (ai *BuildInAiRuntime) StopAi(uuid string) error {
	return nil
}
func (ai *BuildInAiRuntime) Stop() {

}
func (ai *BuildInAiRuntime) Infer([][]float64) [][]float64 {

	return [][]float64{}
}
