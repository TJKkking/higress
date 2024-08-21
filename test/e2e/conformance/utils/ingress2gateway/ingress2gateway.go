package ingress2gateway

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/alibaba/higress/client/pkg/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ToGatewayResource(filepath string) ([]client.Object, error) {
	return ConvertYAMLToClientObjects(ReadFromFile(filepath))
}

func ReadFromFile(filepath string) []byte {
	// basePath := "./conformance/"
	// iPath := basePath + filepath

	cmd := exec.Command("./conformance/utils/ingress2gateway/ingress2gateway", "print",
		"--providers", "higress", "--namespace", "higress-conformance-infra", "--input-file", filepath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}
	return output
}

func ConvertYAMLToClientObjects(yamlData []byte) ([]client.Object, error) {
	// 创建一个解码器
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()

	// 用于存储转换后的对象
	var objects []client.Object

	// fmt.Println("yamlData: ", string(yamlData))

	// 分割多对象的 YAML 数据
	yamlDocs := bytes.Split(yamlData, []byte("\n---\n"))

	for _, doc := range yamlDocs {
		if len(doc) == 0 {
			continue
		}

		obj := &unstructured.Unstructured{}
		_, _, err := decoder.Decode(doc, nil, obj)
		if err != nil {
			return nil, err
		}

		// 将解码后的对象转换为 client.Object
		objects = append(objects, obj)
	}

	return objects, nil
}
