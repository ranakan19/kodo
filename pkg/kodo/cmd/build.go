package cmd

import (
	"context"

	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	buildv1api "github.com/openshift/api/build/v1"
	imagev1api "github.com/openshift/api/image/v1"

	buildv1clientapi "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	imagev1clientapi "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
)

func createTypeMeta(kind string, APIVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: APIVersion,
	}
}

func createObjectType(name string, namespace string) metav1.ObjectMeta {
	if namespace != "" {
		return metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		}
	} else {
		return metav1.ObjectMeta{
			Name: name,
		}
	}
}

func createBuildSpec(uri string) buildv1api.BuildConfigSpec {
	return buildv1api.BuildConfigSpec{
		CommonSpec: buildv1api.CommonSpec{
			Source: buildv1api.BuildSource{
				Type: buildv1api.BuildSourceType("Git"),
				Git: &buildv1api.GitBuildSource{
					URI: uri,
				},
			},
			Strategy: buildv1api.BuildStrategy{
				Type: buildv1api.BuildStrategyType("Docker"),
			},
			Output: buildv1api.BuildOutput{
				To: &corev1.ObjectReference{
					Kind: "ImageStreamTag",
					Name: "my-ruby-image:latest",
				},
			},
		},
	}
}

func createBuildConfig(sourceUrl string) buildv1api.BuildConfig {
	return buildv1api.BuildConfig{
		TypeMeta:   createTypeMeta("BuildConfig", "build.openshift.io/v1"),
		ObjectMeta: createObjectType("my-app-docker-build", ""),
		Spec:       createBuildSpec(sourceUrl),
	}
}

func createImageStream() imagev1api.ImageStream {
	return imagev1api.ImageStream{
		TypeMeta:   createTypeMeta("ImageStream", "image.openshift.io/v1"),
		ObjectMeta: createObjectType("my-ruby-image", "test-kanika"),
	}
}

func newImageStreamClient() *imagev1clientapi.ImageV1Client {
	config := rest.Config{
		Host:        Host,
		BearerToken: Bearertoken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	myClientSet, _ := imagev1clientapi.NewForConfig(&config)
	return myClientSet
}
func newBuildConfigClient() *buildv1clientapi.BuildV1Client {
	config := rest.Config{
		Host:        Host,
		BearerToken: Bearertoken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	myClientSet, _ := buildv1clientapi.NewForConfig(&config)
	return myClientSet
}

//Build image from dockerfile at github source
func Build() error {
	buildclient := newBuildConfigClient()
	buildconfig := createBuildConfig(Source)

	imagestreamclient := newImageStreamClient()
	imagestream := createImageStream()

	_, imgerr := imagestreamclient.ImageStreams(Namespace).Create(context.TODO(), &imagestream, metav1.CreateOptions{})
	_, builderr := buildclient.BuildConfigs(Namespace).Create(context.TODO(), &buildconfig, metav1.CreateOptions{})

	if imgerr != nil {
		return imgerr
	}
	if builderr != nil {
		return builderr
	}
	return nil
}
