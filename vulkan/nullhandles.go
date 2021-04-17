// NULL HANDLES

package vulkan

//import "github.com/vulkan-go/vulkan"
import "unsafe"

/*
#include <vulkan/vulkan.h>

VkPipeline nilPipeline =  (VkPipeline) { 0 };
VkPipelineCache nilPipelineCache =  (VkPipelineCache) { VK_NULL_HANDLE };
*/
import "C"

var NilPipeline = (*Pipeline)(unsafe.Pointer(&C.nilPipeline));
var NilPipelineCache = (*PipelineCache)(unsafe.Pointer(&C.nilPipelineCache));
