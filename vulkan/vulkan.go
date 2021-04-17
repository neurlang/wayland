package vulkan

import "unsafe"

/*

#cgo pkg-config: vulkan
#cgo LDFLAGS: -lvulkan


#include <vulkan/vulkan.h>
#include <vulkan/vulkan_wayland.h>


//////////////////////


VkResult callVkCreateInstance(
    const VkInstanceCreateInfo*                 pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkInstance*                                 pInstance) {
    return vkCreateInstance(pCreateInfo, pAllocator, pInstance);
}

void callVkDestroyInstance(
    VkInstance                                  instance,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyInstance(instance, pAllocator);
}

VkResult callVkEnumeratePhysicalDevices(
    VkInstance                                  instance,
    uint32_t*                                   pPhysicalDeviceCount,
    VkPhysicalDevice*                           pPhysicalDevices) {
    return vkEnumeratePhysicalDevices(instance, pPhysicalDeviceCount, pPhysicalDevices);
}

void callVkGetPhysicalDeviceFeatures(
    VkPhysicalDevice                            physicalDevice,
    VkPhysicalDeviceFeatures*                   pFeatures) {
    vkGetPhysicalDeviceFeatures(physicalDevice, pFeatures);
}

void callVkGetPhysicalDeviceFormatProperties(
    VkPhysicalDevice                            physicalDevice,
    VkFormat                                    format,
    VkFormatProperties*                         pFormatProperties) {
    vkGetPhysicalDeviceFormatProperties(physicalDevice, format, pFormatProperties);
}

VkResult callVkGetPhysicalDeviceImageFormatProperties(
    VkPhysicalDevice                            physicalDevice,
    VkFormat                                    format,
    VkImageType                                 type,
    VkImageTiling                               tiling,
    VkImageUsageFlags                           usage,
    VkImageCreateFlags                          flags,
    VkImageFormatProperties*                    pImageFormatProperties) {
    return vkGetPhysicalDeviceImageFormatProperties(physicalDevice, format, type,
            tiling, usage, flags, pImageFormatProperties);
}

void callVkGetPhysicalDeviceProperties(
    VkPhysicalDevice                            physicalDevice,
    VkPhysicalDeviceProperties*                 pProperties) {
    vkGetPhysicalDeviceProperties(physicalDevice, pProperties);
}

void callVkGetPhysicalDeviceQueueFamilyProperties(
    VkPhysicalDevice                            physicalDevice,
    uint32_t*                                   pQueueFamilyPropertyCount,
    VkQueueFamilyProperties*                    pQueueFamilyProperties) {
    vkGetPhysicalDeviceQueueFamilyProperties(physicalDevice,
            pQueueFamilyPropertyCount, pQueueFamilyProperties);
}

void callVkGetPhysicalDeviceMemoryProperties(
    VkPhysicalDevice                            physicalDevice,
    VkPhysicalDeviceMemoryProperties*           pMemoryProperties) {
    vkGetPhysicalDeviceMemoryProperties(physicalDevice, pMemoryProperties);
}

VkResult callVkCreateDevice(
    VkPhysicalDevice                            physicalDevice,
    const VkDeviceCreateInfo*                   pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDevice*                                   pDevice) {
    return vkCreateDevice(physicalDevice, pCreateInfo, pAllocator, pDevice);
}

void callVkDestroyDevice(
    VkDevice                                    device,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyDevice(device, pAllocator);
}

VkResult callVkEnumerateInstanceExtensionProperties(
    const char*                                 pLayerName,
    uint32_t*                                   pPropertyCount,
    VkExtensionProperties*                      pProperties) {
    return vkEnumerateInstanceExtensionProperties(pLayerName, pPropertyCount, pProperties);
}

VkResult callVkEnumerateDeviceExtensionProperties(
    VkPhysicalDevice                            physicalDevice,
    const char*                                 pLayerName,
    uint32_t*                                   pPropertyCount,
    VkExtensionProperties*                      pProperties) {
    return vkEnumerateDeviceExtensionProperties(physicalDevice, pLayerName,
            pPropertyCount, pProperties);
}

VkResult callVkEnumerateInstanceLayerProperties(
    uint32_t*                                   pPropertyCount,
    VkLayerProperties*                          pProperties) {
    return vkEnumerateInstanceLayerProperties(pPropertyCount, pProperties);
}

VkResult callVkEnumerateDeviceLayerProperties(
    VkPhysicalDevice                            physicalDevice,
    uint32_t*                                   pPropertyCount,
    VkLayerProperties*                          pProperties) {
    return vkEnumerateDeviceLayerProperties(physicalDevice, pPropertyCount, pProperties);
}

void callVkGetDeviceQueue(
    VkDevice                                    device,
    uint32_t                                    queueFamilyIndex,
    uint32_t                                    queueIndex,
    VkQueue*                                    pQueue) {
    vkGetDeviceQueue(device, queueFamilyIndex, queueIndex, pQueue);
}

VkResult callVkQueueSubmit(
    VkQueue                                     queue,
    uint32_t                                    submitCount,
    const VkSubmitInfo*                         pSubmits,
    VkFence                                     fence) {
    return vkQueueSubmit(queue, submitCount, pSubmits, fence);
}

VkResult callVkQueueWaitIdle(
    VkQueue                                     queue) {
    return vkQueueWaitIdle(queue);
}

VkResult callVkDeviceWaitIdle(
    VkDevice                                    device) {
    return vkDeviceWaitIdle(device);
}

VkResult callVkAllocateMemory(
    VkDevice                                    device,
    const VkMemoryAllocateInfo*                 pAllocateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDeviceMemory*                             pMemory) {
    return vkAllocateMemory(device, pAllocateInfo, pAllocator, pMemory);
}

void callVkFreeMemory(
    VkDevice                                    device,
    VkDeviceMemory                              memory,
    const VkAllocationCallbacks*                pAllocator) {
    vkFreeMemory(device, memory, pAllocator);
}

VkResult callVkMapMemory(
    VkDevice                                    device,
    VkDeviceMemory                              memory,
    VkDeviceSize                                offset,
    VkDeviceSize                                size,
    VkMemoryMapFlags                            flags,
    void**                                      ppData) {
    return vkMapMemory(device, memory, offset, size, flags, ppData);
}

void callVkUnmapMemory(
    VkDevice                                    device,
    VkDeviceMemory                              memory) {
    vkUnmapMemory(device, memory);
}

VkResult callVkFlushMappedMemoryRanges(
    VkDevice                                    device,
    uint32_t                                    memoryRangeCount,
    const VkMappedMemoryRange*                  pMemoryRanges) {
    return vkFlushMappedMemoryRanges(device, memoryRangeCount, pMemoryRanges);
}

VkResult callVkInvalidateMappedMemoryRanges(
    VkDevice                                    device,
    uint32_t                                    memoryRangeCount,
    const VkMappedMemoryRange*                  pMemoryRanges) {
    return vkInvalidateMappedMemoryRanges(device, memoryRangeCount, pMemoryRanges);
}

void callVkGetDeviceMemoryCommitment(
    VkDevice                                    device,
    VkDeviceMemory                              memory,
    VkDeviceSize*                               pCommittedMemoryInBytes) {
    vkGetDeviceMemoryCommitment(device, memory, pCommittedMemoryInBytes);
}

VkResult callVkBindBufferMemory(
    VkDevice                                    device,
    VkBuffer                                    buffer,
    VkDeviceMemory                              memory,
    VkDeviceSize                                memoryOffset) {
    return vkBindBufferMemory(device, buffer, memory, memoryOffset);
}

VkResult callVkBindImageMemory(
    VkDevice                                    device,
    VkImage                                     image,
    VkDeviceMemory                              memory,
    VkDeviceSize                                memoryOffset) {
    return vkBindImageMemory(device, image, memory, memoryOffset);
}

void callVkGetBufferMemoryRequirements(
    VkDevice                                    device,
    VkBuffer                                    buffer,
    VkMemoryRequirements*                       pMemoryRequirements) {
    vkGetBufferMemoryRequirements(device, buffer, pMemoryRequirements);
}

void callVkGetImageMemoryRequirements(
    VkDevice                                    device,
    VkImage                                     image,
    VkMemoryRequirements*                       pMemoryRequirements) {
    vkGetImageMemoryRequirements(device, image, pMemoryRequirements);
}

void callVkGetImageSparseMemoryRequirements(
    VkDevice                                    device,
    VkImage                                     image,
    uint32_t*                                   pSparseMemoryRequirementCount,
    VkSparseImageMemoryRequirements*            pSparseMemoryRequirements) {
    vkGetImageSparseMemoryRequirements(device, image, pSparseMemoryRequirementCount,
                                           pSparseMemoryRequirements);
}

void callVkGetPhysicalDeviceSparseImageFormatProperties(
    VkPhysicalDevice                            physicalDevice,
    VkFormat                                    format,
    VkImageType                                 type,
    VkSampleCountFlagBits                       samples,
    VkImageUsageFlags                           usage,
    VkImageTiling                               tiling,
    uint32_t*                                   pPropertyCount,
    VkSparseImageFormatProperties*              pProperties) {
    vkGetPhysicalDeviceSparseImageFormatProperties(physicalDevice, format,
            type, samples, usage, tiling, pPropertyCount, pProperties);
}

VkResult callVkQueueBindSparse(
    VkQueue                                     queue,
    uint32_t                                    bindInfoCount,
    const VkBindSparseInfo*                     pBindInfo,
    VkFence                                     fence) {
    return vkQueueBindSparse(queue, bindInfoCount, pBindInfo, fence);
}

VkResult callVkCreateFence(
    VkDevice                                    device,
    const VkFenceCreateInfo*                    pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkFence*                                    pFence) {
    return vkCreateFence(device, pCreateInfo, pAllocator, pFence);
}

void callVkDestroyFence(
    VkDevice                                    device,
    VkFence                                     fence,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyFence(device, fence, pAllocator);
}

VkResult callVkResetFences(
    VkDevice                                    device,
    uint32_t                                    fenceCount,
    const VkFence*                              pFences) {
    return vkResetFences(device, fenceCount, pFences);
}

VkResult callVkGetFenceStatus(
    VkDevice                                    device,
    VkFence                                     fence) {
    return vkGetFenceStatus(device, fence);
}

VkResult callVkWaitForFences(
    VkDevice                                    device,
    uint32_t                                    fenceCount,
    const VkFence*                              pFences,
    VkBool32                                    waitAll,
    uint64_t                                    timeout) {
    return vkWaitForFences(device, fenceCount, pFences, waitAll, timeout);
}

VkResult callVkCreateSemaphore(
    VkDevice                                    device,
    const VkSemaphoreCreateInfo*                pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSemaphore*                                pSemaphore) {
    return vkCreateSemaphore(device, pCreateInfo, pAllocator, pSemaphore);
}

void callVkDestroySemaphore(
    VkDevice                                    device,
    VkSemaphore                                 semaphore,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroySemaphore(device, semaphore, pAllocator);
}

VkResult callVkCreateEvent(
    VkDevice                                    device,
    const VkEventCreateInfo*                    pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkEvent*                                    pEvent) {
    return vkCreateEvent(device, pCreateInfo, pAllocator, pEvent);
}

void callVkDestroyEvent(
    VkDevice                                    device,
    VkEvent                                     event,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyEvent(device, event, pAllocator);
}

VkResult callVkGetEventStatus(
    VkDevice                                    device,
    VkEvent                                     event) {
    return vkGetEventStatus(device, event);
}

VkResult callVkSetEvent(
    VkDevice                                    device,
    VkEvent                                     event) {
    return vkSetEvent(device, event);
}

VkResult callVkResetEvent(
    VkDevice                                    device,
    VkEvent                                     event) {
    return vkResetEvent(device, event);
}

VkResult callVkCreateQueryPool(
    VkDevice                                    device,
    const VkQueryPoolCreateInfo*                pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkQueryPool*                                pQueryPool) {
    return vkCreateQueryPool(device, pCreateInfo, pAllocator, pQueryPool);
}

void callVkDestroyQueryPool(
    VkDevice                                    device,
    VkQueryPool                                 queryPool,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyQueryPool(device, queryPool, pAllocator);
}

VkResult callVkGetQueryPoolResults(
    VkDevice                                    device,
    VkQueryPool                                 queryPool,
    uint32_t                                    firstQuery,
    uint32_t                                    queryCount,
    size_t                                      dataSize,
    void*                                       pData,
    VkDeviceSize                                stride,
    VkQueryResultFlags                          flags) {
    return vkGetQueryPoolResults(device, queryPool, firstQuery, queryCount,
                                     dataSize, pData, stride, flags);
}

VkResult callVkCreateBuffer(
    VkDevice                                    device,
    const VkBufferCreateInfo*                   pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkBuffer*                                   pBuffer) {
    return vkCreateBuffer(device, pCreateInfo, pAllocator, pBuffer);
}

void callVkDestroyBuffer(
    VkDevice                                    device,
    VkBuffer                                    buffer,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyBuffer(device, buffer, pAllocator);
}

VkResult callVkCreateBufferView(
    VkDevice                                    device,
    const VkBufferViewCreateInfo*               pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkBufferView*                               pView) {
    return vkCreateBufferView(device, pCreateInfo, pAllocator, pView);
}

void callVkDestroyBufferView(
    VkDevice                                    device,
    VkBufferView                                bufferView,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyBufferView(device, bufferView, pAllocator);
}

VkResult callVkCreateImage(
    VkDevice                                    device,
    const VkImageCreateInfo*                    pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkImage*                                    pImage) {
    return vkCreateImage(device, pCreateInfo, pAllocator, pImage);
}

void callVkDestroyImage(
    VkDevice                                    device,
    VkImage                                     image,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyImage(device, image, pAllocator);
}

void callVkGetImageSubresourceLayout(
    VkDevice                                    device,
    VkImage                                     image,
    const VkImageSubresource*                   pSubresource,
    VkSubresourceLayout*                        pLayout) {
    vkGetImageSubresourceLayout(device, image, pSubresource, pLayout);
}

VkResult callVkCreateImageView(
    VkDevice                                    device,
    const VkImageViewCreateInfo*                pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkImageView*                                pView) {
    return vkCreateImageView(device, pCreateInfo, pAllocator, pView);
}

void callVkDestroyImageView(
    VkDevice                                    device,
    VkImageView                                 imageView,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyImageView(device, imageView, pAllocator);
}

VkResult callVkCreateShaderModule(
    VkDevice                                    device,
    const VkShaderModuleCreateInfo*             pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkShaderModule*                             pShaderModule) {
    return vkCreateShaderModule(device, pCreateInfo, pAllocator, pShaderModule);
}

void callVkDestroyShaderModule(
    VkDevice                                    device,
    VkShaderModule                              shaderModule,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyShaderModule(device, shaderModule, pAllocator);
}

VkResult callVkCreatePipelineCache(
    VkDevice                                    device,
    const VkPipelineCacheCreateInfo*            pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkPipelineCache*                            pPipelineCache) {
    return vkCreatePipelineCache(device, pCreateInfo, pAllocator, pPipelineCache);
}

void callVkDestroyPipelineCache(
    VkDevice                                    device,
    VkPipelineCache                             pipelineCache,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyPipelineCache(device, pipelineCache, pAllocator);
}

VkResult callVkGetPipelineCacheData(
    VkDevice                                    device,
    VkPipelineCache                             pipelineCache,
    size_t*                                     pDataSize,
    void*                                       pData) {
    return vkGetPipelineCacheData(device, pipelineCache, pDataSize, pData);
}

VkResult callVkMergePipelineCaches(
    VkDevice                                    device,
    VkPipelineCache                             dstCache,
    uint32_t                                    srcCacheCount,
    const VkPipelineCache*                      pSrcCaches) {
    return vkMergePipelineCaches(device, dstCache, srcCacheCount, pSrcCaches);
}

VkResult callVkCreateGraphicsPipelines(
    VkDevice                                    device,
    VkPipelineCache                             pipelineCache,
    uint32_t                                    createInfoCount,
    const VkGraphicsPipelineCreateInfo*         pCreateInfos,
    const VkAllocationCallbacks*                pAllocator,
    VkPipeline*                                 pPipelines) {
    return vkCreateGraphicsPipelines(device, pipelineCache, createInfoCount,
                                         pCreateInfos, pAllocator, pPipelines);
}

VkResult callVkCreateComputePipelines(
    VkDevice                                    device,
    VkPipelineCache                             pipelineCache,
    uint32_t                                    createInfoCount,
    const VkComputePipelineCreateInfo*          pCreateInfos,
    const VkAllocationCallbacks*                pAllocator,
    VkPipeline*                                 pPipelines) {
    return vkCreateComputePipelines(device, pipelineCache, createInfoCount,
                                        pCreateInfos, pAllocator, pPipelines);
}

void callVkDestroyPipeline(
    VkDevice                                    device,
    VkPipeline                                  pipeline,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyPipeline(device, pipeline, pAllocator);
}

VkResult callVkCreatePipelineLayout(
    VkDevice                                    device,
    const VkPipelineLayoutCreateInfo*           pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkPipelineLayout*                           pPipelineLayout) {
    return vkCreatePipelineLayout(device, pCreateInfo, pAllocator, pPipelineLayout);
}

void callVkDestroyPipelineLayout(
    VkDevice                                    device,
    VkPipelineLayout                            pipelineLayout,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyPipelineLayout(device, pipelineLayout, pAllocator);
}

VkResult callVkCreateSampler(
    VkDevice                                    device,
    const VkSamplerCreateInfo*                  pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSampler*                                  pSampler) {
    return vkCreateSampler(device, pCreateInfo, pAllocator, pSampler);
}

void callVkDestroySampler(
    VkDevice                                    device,
    VkSampler                                   sampler,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroySampler(device, sampler, pAllocator);
}

VkResult callVkCreateDescriptorSetLayout(
    VkDevice                                    device,
    const VkDescriptorSetLayoutCreateInfo*      pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDescriptorSetLayout*                      pSetLayout) {
    return vkCreateDescriptorSetLayout(device, pCreateInfo, pAllocator, pSetLayout);
}

void callVkDestroyDescriptorSetLayout(
    VkDevice                                    device,
    VkDescriptorSetLayout                       descriptorSetLayout,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyDescriptorSetLayout(device, descriptorSetLayout, pAllocator);
}

VkResult callVkCreateDescriptorPool(
    VkDevice                                    device,
    const VkDescriptorPoolCreateInfo*           pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDescriptorPool*                           pDescriptorPool) {
    return vkCreateDescriptorPool(device, pCreateInfo, pAllocator, pDescriptorPool);
}

void callVkDestroyDescriptorPool(
    VkDevice                                    device,
    VkDescriptorPool                            descriptorPool,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyDescriptorPool(device, descriptorPool, pAllocator);
}

VkResult callVkResetDescriptorPool(
    VkDevice                                    device,
    VkDescriptorPool                            descriptorPool,
    VkDescriptorPoolResetFlags                  flags) {
    return vkResetDescriptorPool(device, descriptorPool, flags);
}

VkResult callVkAllocateDescriptorSets(
    VkDevice                                    device,
    const VkDescriptorSetAllocateInfo*          pAllocateInfo,
    VkDescriptorSet*                            pDescriptorSets) {
    return vkAllocateDescriptorSets(device, pAllocateInfo, pDescriptorSets);
}

VkResult callVkFreeDescriptorSets(
    VkDevice                                    device,
    VkDescriptorPool                            descriptorPool,
    uint32_t                                    descriptorSetCount,
    const VkDescriptorSet*                      pDescriptorSets) {
    return vkFreeDescriptorSets(device, descriptorPool, descriptorSetCount, pDescriptorSets);
}

void callVkUpdateDescriptorSets(
    VkDevice                                    device,
    uint32_t                                    descriptorWriteCount,
    const VkWriteDescriptorSet*                 pDescriptorWrites,
    uint32_t                                    descriptorCopyCount,
    const VkCopyDescriptorSet*                  pDescriptorCopies) {
    vkUpdateDescriptorSets(device, descriptorWriteCount, pDescriptorWrites,
                               descriptorCopyCount, pDescriptorCopies);
}

VkResult callVkCreateFramebuffer(
    VkDevice                                    device,
    const VkFramebufferCreateInfo*              pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkFramebuffer*                              pFramebuffer) {
    return vkCreateFramebuffer(device, pCreateInfo, pAllocator, pFramebuffer);
}

void callVkDestroyFramebuffer(
    VkDevice                                    device,
    VkFramebuffer                               framebuffer,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyFramebuffer(device, framebuffer, pAllocator);
}

VkResult callVkCreateRenderPass(
    VkDevice                                    device,
    const VkRenderPassCreateInfo*               pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkRenderPass*                               pRenderPass) {
    return vkCreateRenderPass(device, pCreateInfo, pAllocator, pRenderPass);
}

void callVkDestroyRenderPass(
    VkDevice                                    device,
    VkRenderPass                                renderPass,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyRenderPass(device, renderPass, pAllocator);
}

void callVkGetRenderAreaGranularity(
    VkDevice                                    device,
    VkRenderPass                                renderPass,
    VkExtent2D*                                 pGranularity) {
    vkGetRenderAreaGranularity(device, renderPass, pGranularity);
}

VkResult callVkCreateCommandPool(
    VkDevice                                    device,
    const VkCommandPoolCreateInfo*              pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkCommandPool*                              pCommandPool) {
    return vkCreateCommandPool(device, pCreateInfo, pAllocator, pCommandPool);
}

void callVkDestroyCommandPool(
    VkDevice                                    device,
    VkCommandPool                               commandPool,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroyCommandPool(device, commandPool, pAllocator);
}

VkResult callVkResetCommandPool(
    VkDevice                                    device,
    VkCommandPool                               commandPool,
    VkCommandPoolResetFlags                     flags) {
    return vkResetCommandPool(device, commandPool, flags);
}

VkResult callVkAllocateCommandBuffers(
    VkDevice                                    device,
    const VkCommandBufferAllocateInfo*          pAllocateInfo,
    VkCommandBuffer*                            pCommandBuffers) {
    return vkAllocateCommandBuffers(device, pAllocateInfo, pCommandBuffers);
}

void callVkFreeCommandBuffers(
    VkDevice                                    device,
    VkCommandPool                               commandPool,
    uint32_t                                    commandBufferCount,
    const VkCommandBuffer*                      pCommandBuffers) {
    vkFreeCommandBuffers(device, commandPool, commandBufferCount, pCommandBuffers);
}

VkResult callVkBeginCommandBuffer(
    VkCommandBuffer                             commandBuffer,
    const VkCommandBufferBeginInfo*             pBeginInfo) {
    return vkBeginCommandBuffer(commandBuffer, pBeginInfo);
}

VkResult callVkEndCommandBuffer(
    VkCommandBuffer                             commandBuffer) {
    return vkEndCommandBuffer(commandBuffer);
}

VkResult callVkResetCommandBuffer(
    VkCommandBuffer                             commandBuffer,
    VkCommandBufferResetFlags                   flags) {
    return vkResetCommandBuffer(commandBuffer, flags);
}

void callVkCmdBindPipeline(
    VkCommandBuffer                             commandBuffer,
    VkPipelineBindPoint                         pipelineBindPoint,
    VkPipeline                                  pipeline) {
    vkCmdBindPipeline(commandBuffer, pipelineBindPoint, pipeline);
}

void callVkCmdSetViewport(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    firstViewport,
    uint32_t                                    viewportCount,
    const VkViewport*                           pViewports) {
    vkCmdSetViewport(commandBuffer, firstViewport, viewportCount, pViewports);
}

void callVkCmdSetScissor(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    firstScissor,
    uint32_t                                    scissorCount,
    const VkRect2D*                             pScissors) {
    vkCmdSetScissor(commandBuffer, firstScissor, scissorCount, pScissors);
}

void callVkCmdSetLineWidth(
    VkCommandBuffer                             commandBuffer,
    float                                       lineWidth) {
    vkCmdSetLineWidth(commandBuffer, lineWidth);
}

void callVkCmdSetDepthBias(
    VkCommandBuffer                             commandBuffer,
    float                                       depthBiasConstantFactor,
    float                                       depthBiasClamp,
    float                                       depthBiasSlopeFactor) {
    vkCmdSetDepthBias(commandBuffer, depthBiasConstantFactor,
                          depthBiasClamp, depthBiasSlopeFactor);
}

void callVkCmdSetBlendConstants(
    VkCommandBuffer                             commandBuffer,
    const float                                 blendConstants[4]) {
    vkCmdSetBlendConstants(commandBuffer, blendConstants);
}

void callVkCmdSetDepthBounds(
    VkCommandBuffer                             commandBuffer,
    float                                       minDepthBounds,
    float                                       maxDepthBounds) {
    vkCmdSetDepthBounds(commandBuffer, minDepthBounds, maxDepthBounds);
}

void callVkCmdSetStencilCompareMask(
    VkCommandBuffer                             commandBuffer,
    VkStencilFaceFlags                          faceMask,
    uint32_t                                    compareMask) {
    vkCmdSetStencilCompareMask(commandBuffer, faceMask, compareMask);
}

void callVkCmdSetStencilWriteMask(
    VkCommandBuffer                             commandBuffer,
    VkStencilFaceFlags                          faceMask,
    uint32_t                                    writeMask) {
    vkCmdSetStencilWriteMask(commandBuffer, faceMask, writeMask);
}

void callVkCmdSetStencilReference(
    VkCommandBuffer                             commandBuffer,
    VkStencilFaceFlags                          faceMask,
    uint32_t                                    reference) {
    vkCmdSetStencilReference(commandBuffer, faceMask, reference);
}

void callVkCmdBindDescriptorSets(
    VkCommandBuffer                             commandBuffer,
    VkPipelineBindPoint                         pipelineBindPoint,
    VkPipelineLayout                            layout,
    uint32_t                                    firstSet,
    uint32_t                                    descriptorSetCount,
    const VkDescriptorSet*                      pDescriptorSets,
    uint32_t                                    dynamicOffsetCount,
    const uint32_t*                             pDynamicOffsets) {
    vkCmdBindDescriptorSets(commandBuffer, pipelineBindPoint, layout,
                                firstSet, descriptorSetCount, pDescriptorSets,
                                dynamicOffsetCount, pDynamicOffsets);
}

void callVkCmdBindIndexBuffer(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    buffer,
    VkDeviceSize                                offset,
    VkIndexType                                 indexType) {
    vkCmdBindIndexBuffer(commandBuffer, buffer, offset, indexType);
}

void callVkCmdBindVertexBuffers(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    firstBinding,
    uint32_t                                    bindingCount,
    const VkBuffer*                             pBuffers,
    const VkDeviceSize*                         pOffsets) {
    vkCmdBindVertexBuffers(commandBuffer, firstBinding, bindingCount, pBuffers, pOffsets);
}

void callVkCmdDraw(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    vertexCount,
    uint32_t                                    instanceCount,
    uint32_t                                    firstVertex,
    uint32_t                                    firstInstance) {
    vkCmdDraw(commandBuffer, vertexCount, instanceCount, firstVertex, firstInstance);
}

void callVkCmdDrawIndexed(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    indexCount,
    uint32_t                                    instanceCount,
    uint32_t                                    firstIndex,
    int32_t                                     vertexOffset,
    uint32_t                                    firstInstance) {
    vkCmdDrawIndexed(commandBuffer, indexCount, instanceCount,
                         firstIndex, vertexOffset, firstInstance);
}

void callVkCmdDrawIndirect(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    buffer,
    VkDeviceSize                                offset,
    uint32_t                                    drawCount,
    uint32_t                                    stride) {
    vkCmdDrawIndirect(commandBuffer, buffer, offset, drawCount, stride);
}

void callVkCmdDrawIndexedIndirect(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    buffer,
    VkDeviceSize                                offset,
    uint32_t                                    drawCount,
    uint32_t                                    stride) {
    vkCmdDrawIndexedIndirect(commandBuffer, buffer, offset, drawCount, stride);
}

void callVkCmdDispatch(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    x,
    uint32_t                                    y,
    uint32_t                                    z) {
    vkCmdDispatch(commandBuffer, x, y, z);
}

void callVkCmdDispatchIndirect(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    buffer,
    VkDeviceSize                                offset) {
    vkCmdDispatchIndirect(commandBuffer, buffer, offset);
}

void callVkCmdCopyBuffer(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    srcBuffer,
    VkBuffer                                    dstBuffer,
    uint32_t                                    regionCount,
    const VkBufferCopy*                         pRegions) {
    vkCmdCopyBuffer(commandBuffer, srcBuffer, dstBuffer, regionCount, pRegions);
}

void callVkCmdCopyImage(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     srcImage,
    VkImageLayout                               srcImageLayout,
    VkImage                                     dstImage,
    VkImageLayout                               dstImageLayout,
    uint32_t                                    regionCount,
    const VkImageCopy*                          pRegions) {
    vkCmdCopyImage(commandBuffer, srcImage, srcImageLayout,
                       dstImage, dstImageLayout, regionCount, pRegions);
}

void callVkCmdBlitImage(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     srcImage,
    VkImageLayout                               srcImageLayout,
    VkImage                                     dstImage,
    VkImageLayout                               dstImageLayout,
    uint32_t                                    regionCount,
    const VkImageBlit*                          pRegions,
    VkFilter                                    filter) {
    vkCmdBlitImage(commandBuffer, srcImage, srcImageLayout,
                       dstImage, dstImageLayout, regionCount, pRegions, filter);
}

void callVkCmdCopyBufferToImage(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    srcBuffer,
    VkImage                                     dstImage,
    VkImageLayout                               dstImageLayout,
    uint32_t                                    regionCount,
    const VkBufferImageCopy*                    pRegions) {
    vkCmdCopyBufferToImage(commandBuffer, srcBuffer,
                               dstImage, dstImageLayout, regionCount, pRegions);
}

void callVkCmdCopyImageToBuffer(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     srcImage,
    VkImageLayout                               srcImageLayout,
    VkBuffer                                    dstBuffer,
    uint32_t                                    regionCount,
    const VkBufferImageCopy*                    pRegions) {
    vkCmdCopyImageToBuffer(commandBuffer, srcImage, srcImageLayout,
                               dstBuffer, regionCount, pRegions);
}

void callVkCmdUpdateBuffer(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    dstBuffer,
    VkDeviceSize                                dstOffset,
    VkDeviceSize                                dataSize,
    const uint32_t*                             pData) {
    vkCmdUpdateBuffer(commandBuffer, dstBuffer, dstOffset, dataSize, pData);
}

void callVkCmdFillBuffer(
    VkCommandBuffer                             commandBuffer,
    VkBuffer                                    dstBuffer,
    VkDeviceSize                                dstOffset,
    VkDeviceSize                                size,
    uint32_t                                    data) {
    vkCmdFillBuffer(commandBuffer, dstBuffer, dstOffset, size, data);
}

void callVkCmdClearColorImage(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     image,
    VkImageLayout                               imageLayout,
    const VkClearColorValue*                    pColor,
    uint32_t                                    rangeCount,
    const VkImageSubresourceRange*              pRanges) {
    vkCmdClearColorImage(commandBuffer, image, imageLayout, pColor, rangeCount, pRanges);
}

void callVkCmdClearDepthStencilImage(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     image,
    VkImageLayout                               imageLayout,
    const VkClearDepthStencilValue*             pDepthStencil,
    uint32_t                                    rangeCount,
    const VkImageSubresourceRange*              pRanges) {
    vkCmdClearDepthStencilImage(commandBuffer, image, imageLayout,
                                    pDepthStencil, rangeCount, pRanges);
}

void callVkCmdClearAttachments(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    attachmentCount,
    const VkClearAttachment*                    pAttachments,
    uint32_t                                    rectCount,
    const VkClearRect*                          pRects) {
    vkCmdClearAttachments(commandBuffer, attachmentCount, pAttachments, rectCount, pRects);
}

void callVkCmdResolveImage(
    VkCommandBuffer                             commandBuffer,
    VkImage                                     srcImage,
    VkImageLayout                               srcImageLayout,
    VkImage                                     dstImage,
    VkImageLayout                               dstImageLayout,
    uint32_t                                    regionCount,
    const VkImageResolve*                       pRegions) {
    vkCmdResolveImage(commandBuffer, srcImage, srcImageLayout,
                          dstImage, dstImageLayout, regionCount, pRegions);
}

void callVkCmdSetEvent(
    VkCommandBuffer                             commandBuffer,
    VkEvent                                     event,
    VkPipelineStageFlags                        stageMask) {
    vkCmdSetEvent(commandBuffer, event, stageMask);
}

void callVkCmdResetEvent(
    VkCommandBuffer                             commandBuffer,
    VkEvent                                     event,
    VkPipelineStageFlags                        stageMask) {
    vkCmdResetEvent(commandBuffer, event, stageMask);
}

void callVkCmdWaitEvents(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    eventCount,
    const VkEvent*                              pEvents,
    VkPipelineStageFlags                        srcStageMask,
    VkPipelineStageFlags                        dstStageMask,
    uint32_t                                    memoryBarrierCount,
    const VkMemoryBarrier*                      pMemoryBarriers,
    uint32_t                                    bufferMemoryBarrierCount,
    const VkBufferMemoryBarrier*                pBufferMemoryBarriers,
    uint32_t                                    imageMemoryBarrierCount,
    const VkImageMemoryBarrier*                 pImageMemoryBarriers) {
    vkCmdWaitEvents(commandBuffer, eventCount, pEvents, srcStageMask, dstStageMask,
                        memoryBarrierCount, pMemoryBarriers,
                        bufferMemoryBarrierCount, pBufferMemoryBarriers,
                        imageMemoryBarrierCount, pImageMemoryBarriers);
}

void callVkCmdPipelineBarrier(
    VkCommandBuffer                             commandBuffer,
    VkPipelineStageFlags                        srcStageMask,
    VkPipelineStageFlags                        dstStageMask,
    VkDependencyFlags                           dependencyFlags,
    uint32_t                                    memoryBarrierCount,
    const VkMemoryBarrier*                      pMemoryBarriers,
    uint32_t                                    bufferMemoryBarrierCount,
    const VkBufferMemoryBarrier*                pBufferMemoryBarriers,
    uint32_t                                    imageMemoryBarrierCount,
    const VkImageMemoryBarrier*                 pImageMemoryBarriers) {
    vkCmdPipelineBarrier(commandBuffer, srcStageMask, dstStageMask, dependencyFlags,
                             memoryBarrierCount, pMemoryBarriers,
                             bufferMemoryBarrierCount, pBufferMemoryBarriers,
                             imageMemoryBarrierCount, pImageMemoryBarriers);
}

void callVkCmdBeginQuery(
    VkCommandBuffer                             commandBuffer,
    VkQueryPool                                 queryPool,
    uint32_t                                    query,
    VkQueryControlFlags                         flags) {
    vkCmdBeginQuery(commandBuffer, queryPool, query, flags);
}

void callVkCmdEndQuery(
    VkCommandBuffer                             commandBuffer,
    VkQueryPool                                 queryPool,
    uint32_t                                    query) {
    vkCmdEndQuery(commandBuffer, queryPool, query);
}

void callVkCmdResetQueryPool(
    VkCommandBuffer                             commandBuffer,
    VkQueryPool                                 queryPool,
    uint32_t                                    firstQuery,
    uint32_t                                    queryCount) {
    vkCmdResetQueryPool(commandBuffer, queryPool, firstQuery, queryCount);
}

void callVkCmdWriteTimestamp(
    VkCommandBuffer                             commandBuffer,
    VkPipelineStageFlagBits                     pipelineStage,
    VkQueryPool                                 queryPool,
    uint32_t                                    query) {
    vkCmdWriteTimestamp(commandBuffer, pipelineStage, queryPool, query);
}

void callVkCmdCopyQueryPoolResults(
    VkCommandBuffer                             commandBuffer,
    VkQueryPool                                 queryPool,
    uint32_t                                    firstQuery,
    uint32_t                                    queryCount,
    VkBuffer                                    dstBuffer,
    VkDeviceSize                                dstOffset,
    VkDeviceSize                                stride,
    VkQueryResultFlags                          flags) {
    vkCmdCopyQueryPoolResults(commandBuffer, queryPool, firstQuery, queryCount,
                                  dstBuffer, dstOffset, stride, flags);
}

void callVkCmdPushConstants(
    VkCommandBuffer                             commandBuffer,
    VkPipelineLayout                            layout,
    VkShaderStageFlags                          stageFlags,
    uint32_t                                    offset,
    uint32_t                                    size,
    const void*                                 pValues) {
    vkCmdPushConstants(commandBuffer, layout, stageFlags, offset, size, pValues);
}

void callVkCmdBeginRenderPass(
    VkCommandBuffer                             commandBuffer,
    const VkRenderPassBeginInfo*                pRenderPassBegin,
    VkSubpassContents                           contents) {
    vkCmdBeginRenderPass(commandBuffer, pRenderPassBegin, contents);
}

void callVkCmdNextSubpass(
    VkCommandBuffer                             commandBuffer,
    VkSubpassContents                           contents) {
    vkCmdNextSubpass(commandBuffer, contents);
}

void callVkCmdEndRenderPass(
    VkCommandBuffer                             commandBuffer) {
    vkCmdEndRenderPass(commandBuffer);
}

void callVkCmdExecuteCommands(
    VkCommandBuffer                             commandBuffer,
    uint32_t                                    commandBufferCount,
    const VkCommandBuffer*                      pCommandBuffers) {
    vkCmdExecuteCommands(commandBuffer, commandBufferCount, pCommandBuffers);
}

void callVkDestroySurfaceKHR(
    VkInstance                                  instance,
    VkSurfaceKHR                                surface,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroySurfaceKHR(instance, surface, pAllocator);
}

VkResult callVkGetPhysicalDeviceSurfaceSupportKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex,
    VkSurfaceKHR                                surface,
    VkBool32*                                   pSupported) {
    return vkGetPhysicalDeviceSurfaceSupportKHR(physicalDevice,
            queueFamilyIndex, surface, pSupported);
}

VkResult callVkGetPhysicalDeviceSurfaceCapabilitiesKHR(
    VkPhysicalDevice                            physicalDevice,
    VkSurfaceKHR                                surface,
    VkSurfaceCapabilitiesKHR*                   pSurfaceCapabilities) {
    return vkGetPhysicalDeviceSurfaceCapabilitiesKHR(physicalDevice,
            surface, pSurfaceCapabilities);
}

VkResult callVkGetPhysicalDeviceSurfaceFormatsKHR(
    VkPhysicalDevice                            physicalDevice,
    VkSurfaceKHR                                surface,
    uint32_t*                                   pSurfaceFormatCount,
    VkSurfaceFormatKHR*                         pSurfaceFormats) {
    return vkGetPhysicalDeviceSurfaceFormatsKHR(physicalDevice,
            surface, pSurfaceFormatCount, pSurfaceFormats);
}

VkResult callVkGetPhysicalDeviceSurfacePresentModesKHR(
    VkPhysicalDevice                            physicalDevice,
    VkSurfaceKHR                                surface,
    uint32_t*                                   pPresentModeCount,
    VkPresentModeKHR*                           pPresentModes) {
    return vkGetPhysicalDeviceSurfacePresentModesKHR(physicalDevice,
            surface, pPresentModeCount, pPresentModes);
}

VkResult callVkCreateSwapchainKHR(
    VkDevice                                    device,
    const VkSwapchainCreateInfoKHR*             pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSwapchainKHR*                             pSwapchain) {
    return vkCreateSwapchainKHR(device, pCreateInfo, pAllocator, pSwapchain);
}

void callVkDestroySwapchainKHR(
    VkDevice                                    device,
    VkSwapchainKHR                              swapchain,
    const VkAllocationCallbacks*                pAllocator) {
    vkDestroySwapchainKHR(device, swapchain, pAllocator);
}

VkResult callVkGetSwapchainImagesKHR(
    VkDevice                                    device,
    VkSwapchainKHR                              swapchain,
    uint32_t*                                   pSwapchainImageCount,
    VkImage*                                    pSwapchainImages) {
    return vkGetSwapchainImagesKHR(device, swapchain, pSwapchainImageCount, pSwapchainImages);
}

VkResult callVkAcquireNextImageKHR(
    VkDevice                                    device,
    VkSwapchainKHR                              swapchain,
    uint64_t                                    timeout,
    VkSemaphore                                 semaphore,
    VkFence                                     fence,
    uint32_t*                                   pImageIndex) {
    return vkAcquireNextImageKHR(device, swapchain, timeout, semaphore, fence, pImageIndex);
}

VkResult callVkQueuePresentKHR(
    VkQueue                                     queue,
    const VkPresentInfoKHR*                     pPresentInfo) {
    return vkQueuePresentKHR(queue, pPresentInfo);
}

VkResult callVkGetPhysicalDeviceDisplayPropertiesKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t*                                   pPropertyCount,
    VkDisplayPropertiesKHR*                     pProperties) {
    return vkGetPhysicalDeviceDisplayPropertiesKHR(physicalDevice,
            pPropertyCount, pProperties);
}

VkResult callVkGetPhysicalDeviceDisplayPlanePropertiesKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t*                                   pPropertyCount,
    VkDisplayPlanePropertiesKHR*                pProperties) {
    return vkGetPhysicalDeviceDisplayPlanePropertiesKHR(physicalDevice,
            pPropertyCount, pProperties);
}

VkResult callVkGetDisplayPlaneSupportedDisplaysKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    planeIndex,
    uint32_t*                                   pDisplayCount,
    VkDisplayKHR*                               pDisplays) {
    return vkGetDisplayPlaneSupportedDisplaysKHR(physicalDevice, planeIndex,
            pDisplayCount, pDisplays);
}

VkResult callVkGetDisplayModePropertiesKHR(
    VkPhysicalDevice                            physicalDevice,
    VkDisplayKHR                                display,
    uint32_t*                                   pPropertyCount,
    VkDisplayModePropertiesKHR*                 pProperties) {
    return vkGetDisplayModePropertiesKHR(physicalDevice, display,
            pPropertyCount, pProperties);
}

VkResult callVkCreateDisplayModeKHR(
    VkPhysicalDevice                            physicalDevice,
    VkDisplayKHR                                display,
    const VkDisplayModeCreateInfoKHR*           pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDisplayModeKHR*                           pMode) {
    return vkCreateDisplayModeKHR(physicalDevice, display, pCreateInfo, pAllocator, pMode);
}

VkResult callVkGetDisplayPlaneCapabilitiesKHR(
    VkPhysicalDevice                            physicalDevice,
    VkDisplayModeKHR                            mode,
    uint32_t                                    planeIndex,
    VkDisplayPlaneCapabilitiesKHR*              pCapabilities) {
    return vkGetDisplayPlaneCapabilitiesKHR(physicalDevice, mode, planeIndex, pCapabilities);
}

VkResult callVkCreateDisplayPlaneSurfaceKHR(
    VkInstance                                  instance,
    const VkDisplaySurfaceCreateInfoKHR*        pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateDisplayPlaneSurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}

VkResult callVkCreateSharedSwapchainsKHR(
    VkDevice                                    device,
    uint32_t                                    swapchainCount,
    const VkSwapchainCreateInfoKHR*             pCreateInfos,
    const VkAllocationCallbacks*                pAllocator,
    VkSwapchainKHR*                             pSwapchains) {
    return vkCreateSharedSwapchainsKHR(device, swapchainCount, pCreateInfos,
                                           pAllocator, pSwapchains);
}

#ifdef VK_USE_PLATFORM_XLIB_KHR
VkResult callVkCreateXlibSurfaceKHR(
    VkInstance                                  instance,
    const VkXlibSurfaceCreateInfoKHR*           pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateXlibSurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}

VkBool32 callVkGetPhysicalDeviceXlibPresentationSupportKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex,
    Display*                                    dpy,
    VisualID                                    visualID) {
    return vkGetPhysicalDeviceXlibPresentationSupportKHR(physicalDevice,
            queueFamilyIndex, dpy, visualID);
}
#endif

#ifdef VK_USE_PLATFORM_XCB_KHR
VkResult callVkCreateXcbSurfaceKHR(
    VkInstance                                  instance,
    const VkXcbSurfaceCreateInfoKHR*            pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateXcbSurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}

VkBool32 callVkGetPhysicalDeviceXcbPresentationSupportKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex,
    xcb_connection_t*                           connection,
    xcb_visualid_t                              visual_id) {
    vkGetPhysicalDeviceXcbPresentationSupportKHR(physicalDevice,
            queueFamilyIndex, connection, visual_id);
}
#endif

#ifdef VK_USE_PLATFORM_WAYLAND_KHR
VkResult callVkCreateWaylandSurfaceKHR(
    VkInstance                                  instance,
    const VkWaylandSurfaceCreateInfoKHR*        pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateWaylandSurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}

VkBool32 callVkGetPhysicalDeviceWaylandPresentationSupportKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex,
    struct wl_display*                          display) {
    return vkGetPhysicalDeviceWaylandPresentationSupportKHR(physicalDevice,
            queueFamilyIndex, display);
}
#endif

#ifdef VK_USE_PLATFORM_ANDROID_KHR
VkResult callVkCreateAndroidSurfaceKHR(
    VkInstance                                  instance,
    const VkAndroidSurfaceCreateInfoKHR*        pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateAndroidSurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}
#endif

#ifdef VK_USE_PLATFORM_IOS_MVK
VkResult callVkCreateIOSSurfaceMVK(
    VkInstance                              instance,
    const VkIOSSurfaceCreateInfoMVK*        pCreateInfo,
    const VkAllocationCallbacks*            pAllocator,
    VkSurfaceKHR*                           pSurface) {
    return vkCreateIOSSurfaceMVK(instance, pCreateInfo, pAllocator, pSurface);
}

VkResult callVkActivateMoltenVKLicenseMVK(
    const char*                                 licenseID,
    const char*                                 licenseKey,
    VkBool32                                    acceptLicenseTermsAndConditions) {
    return vkActivateMoltenVKLicenseMVK(licenseID, licenseKey, acceptLicenseTermsAndConditions);
}

VkResult callVkActivateMoltenVKLicensesMVK() {
    return vkActivateMoltenVKLicensesMVK();
}

VkResult callVkGetMoltenVKDeviceConfigurationMVK(
    VkDevice                                    device,
    MVKDeviceConfiguration*                     pConfiguration) {
    return vkGetMoltenVKDeviceConfigurationMVK(device, pConfiguration);
}

VkResult callVkSetMoltenVKDeviceConfigurationMVK(
    VkDevice                                    device,
    MVKDeviceConfiguration*                     pConfiguration) {
    return vkSetMoltenVKDeviceConfigurationMVK(device, pConfiguration);
}

VkResult callVkGetPhysicalDeviceMetalFeaturesMVK(
    VkPhysicalDevice                            physicalDevice,
    MVKPhysicalDeviceMetalFeatures*             pMetalFeatures) {
    return vkGetPhysicalDeviceMetalFeaturesMVK(physicalDevice, pMetalFeatures);
}

VkResult callVkGetSwapchainPerformanceMVK(
    VkDevice                                    device,
    VkSwapchainKHR                              swapchain,
    MVKSwapchainPerformance*                    pSwapchainPerf) {
    return vkGetSwapchainPerformanceMVK(device, swapchain, pSwapchainPerf);
}
#endif

#ifdef VK_USE_PLATFORM_WIN32_KHR
VkResult callVkCreateWin32SurfaceKHR(
    VkInstance                                  instance,
    const VkWin32SurfaceCreateInfoKHR*          pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
    return vkCreateWin32SurfaceKHR(instance, pCreateInfo, pAllocator, pSurface);
}

VkBool32 callVkGetPhysicalDeviceWin32PresentationSupportKHR(
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex) {
    return vkGetPhysicalDeviceWin32PresentationSupportKHR(physicalDevice, queueFamilyIndex);
}
#endif

VkResult callVkCreateDebugReportCallbackEXT(
    VkInstance                                  instance,
    const VkDebugReportCallbackCreateInfoEXT*   pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkDebugReportCallbackEXT*                   pCallback) {

    PFN_vkCreateDebugReportCallbackEXT pfn = (PFN_vkCreateDebugReportCallbackEXT)
            (vkGetInstanceProcAddr(instance, "vkCreateDebugReportCallbackEXT"));
    if (pfn != NULL) {
        return pfn(instance, pCreateInfo, pAllocator, pCallback);
    }
    return VK_NOT_READY;
}

void callVkDestroyDebugReportCallbackEXT(
    VkInstance                                  instance,
    VkDebugReportCallbackEXT                    callback,
    const VkAllocationCallbacks*                pAllocator) {

    PFN_vkDestroyDebugReportCallbackEXT pfn = (PFN_vkDestroyDebugReportCallbackEXT)
            (vkGetInstanceProcAddr(instance, "vkDestroyDebugReportCallbackEXT"));
    if (pfn != NULL) {
        pfn(instance, callback, pAllocator);
    }
}

void callVkDebugReportMessageEXT(
    VkInstance                                  instance,
    VkDebugReportFlagsEXT                       flags,
    VkDebugReportObjectTypeEXT                  objectType,
    uint64_t                                    object,
    size_t                                      location,
    int32_t                                     messageCode,
    const char*                                 pLayerPrefix,
    const char*                                 pMessage) {

    PFN_vkDebugReportMessageEXT pfn = (PFN_vkDebugReportMessageEXT)
                                      (vkGetInstanceProcAddr(instance, "vkDebugReportMessageEXT"));
    if (pfn != NULL) {
        pfn(instance, flags, objectType, object, location,
            messageCode, pLayerPrefix, pMessage);
    }
}


//////////////////////
void callVkGetPhysicalDeviceFeatures2(
    VkPhysicalDevice                            physicalDevice,
    VkPhysicalDeviceFeatures2*                   pFeatures) {
    vkGetPhysicalDeviceFeatures2(physicalDevice, pFeatures);
}

VkResult
callCreateWaylandSurfaceKHR(
    VkInstance                                  instance,
    const VkWaylandSurfaceCreateInfoKHR*        pCreateInfo,
    const VkAllocationCallbacks*                pAllocator,
    VkSurfaceKHR*                               pSurface) {
   PFN_vkCreateWaylandSurfaceKHR pfn =
      (PFN_vkCreateWaylandSurfaceKHR)
      (vkGetInstanceProcAddr(instance, "vkCreateWaylandSurfaceKHR"));
    if (pfn != NULL) {
        return pfn(instance, pCreateInfo, pAllocator, pSurface);
    }
    return VK_NOT_READY;
}
VkBool32 callGetPhysicalDeviceWaylandPresentationSupportKHR(
    VkInstance                                  instance,
    VkPhysicalDevice                            physicalDevice,
    uint32_t                                    queueFamilyIndex,
    struct wl_display*                          display) {
   PFN_vkGetPhysicalDeviceWaylandPresentationSupportKHR pfn =
      (PFN_vkGetPhysicalDeviceWaylandPresentationSupportKHR)
      vkGetInstanceProcAddr(instance, "vkGetPhysicalDeviceWaylandPresentationSupportKHR");
    if (pfn != NULL) {
        return pfn(physicalDevice, queueFamilyIndex, display);
    }
    return VK_NOT_READY;
}
//////////////////////
*/
import "C"

// AcquireNextImage function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkAcquireNextImageKHR
func AcquireNextImage(device Device, swapchain Swapchain, timeout uint64, semaphore Semaphore, fence Fence, pImageIndex *uint32) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cswapchain, _ := *(*C.VkSwapchainKHR)(unsafe.Pointer(&swapchain)), cgoAllocsUnknown
	ctimeout, _ := (C.uint64_t)(timeout), cgoAllocsUnknown
	csemaphore, _ := *(*C.VkSemaphore)(unsafe.Pointer(&semaphore)), cgoAllocsUnknown
	cfence, _ := *(*C.VkFence)(unsafe.Pointer(&fence)), cgoAllocsUnknown
	cpImageIndex, _ := (*C.uint32_t)(unsafe.Pointer(pImageIndex)), cgoAllocsUnknown
	__ret := C.callVkAcquireNextImageKHR(cdevice, cswapchain, ctimeout, csemaphore, cfence, cpImageIndex)
	__v := (Result)(__ret)
	return __v
}

// QueuePresent function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkQueuePresentKHR
func QueuePresent(queue Queue, pPresentInfo *PresentInfo) Result {
	cqueue, _ := *(*C.VkQueue)(unsafe.Pointer(&queue)), cgoAllocsUnknown
	cpPresentInfo, _ := pPresentInfo.PassRef()
	__ret := C.callVkQueuePresentKHR(cqueue, cpPresentInfo)
	__v := (Result)(__ret)
	return __v
}

// QueueWaitIdle function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkQueueWaitIdle.html
func QueueWaitIdle(queue Queue) Result {
	cqueue, _ := *(*C.VkQueue)(unsafe.Pointer(&queue)), cgoAllocsUnknown
	__ret := C.callVkQueueWaitIdle(cqueue)
	__v := (Result)(__ret)
	return __v
}

// WaitForFences function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkWaitForFences.html
func WaitForFences(device Device, fenceCount uint32, pFences []Fence, waitAll Bool32, timeout uint64) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cfenceCount, _ := (C.uint32_t)(fenceCount), cgoAllocsUnknown
	cpFences, _ := (*C.VkFence)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pFences)).Data)), cgoAllocsUnknown
	cwaitAll, _ := (C.VkBool32)(waitAll), cgoAllocsUnknown
	ctimeout, _ := (C.uint64_t)(timeout), cgoAllocsUnknown
	__ret := C.callVkWaitForFences(cdevice, cfenceCount, cpFences, cwaitAll, ctimeout)
	__v := (Result)(__ret)
	return __v
}

// ResetFences function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkResetFences.html
func ResetFences(device Device, fenceCount uint32, pFences []Fence) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cfenceCount, _ := (C.uint32_t)(fenceCount), cgoAllocsUnknown
	cpFences, _ := (*C.VkFence)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pFences)).Data)), cgoAllocsUnknown
	__ret := C.callVkResetFences(cdevice, cfenceCount, cpFences)
	__v := (Result)(__ret)
	return __v
}

// BeginCommandBuffer function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkBeginCommandBuffer.html
func BeginCommandBuffer(commandBuffer CommandBuffer, pBeginInfo *CommandBufferBeginInfo) Result {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cpBeginInfo, _ := pBeginInfo.PassRef()
	__ret := C.callVkBeginCommandBuffer(ccommandBuffer, cpBeginInfo)
	__v := (Result)(__ret)
	return __v
}

// CmdBeginRenderPass function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdBeginRenderPass.html
func CmdBeginRenderPass(commandBuffer CommandBuffer, pRenderPassBegin *RenderPassBeginInfo, contents SubpassContents) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cpRenderPassBegin, _ := pRenderPassBegin.PassRef()
	ccontents, _ := (C.VkSubpassContents)(contents), cgoAllocsUnknown
	C.callVkCmdBeginRenderPass(ccommandBuffer, cpRenderPassBegin, ccontents)
}

// CmdBindVertexBuffers function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdBindVertexBuffers.html
func CmdBindVertexBuffers(commandBuffer CommandBuffer, firstBinding uint32, bindingCount uint32, pBuffers []Buffer, pOffsets []DeviceSize) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cfirstBinding, _ := (C.uint32_t)(firstBinding), cgoAllocsUnknown
	cbindingCount, _ := (C.uint32_t)(bindingCount), cgoAllocsUnknown
	cpBuffers, _ := (*C.VkBuffer)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pBuffers)).Data)), cgoAllocsUnknown
	cpOffsets, _ := (*C.VkDeviceSize)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pOffsets)).Data)), cgoAllocsUnknown
	C.callVkCmdBindVertexBuffers(ccommandBuffer, cfirstBinding, cbindingCount, cpBuffers, cpOffsets)
}

// CmdDraw function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdDraw.html
func CmdDraw(commandBuffer CommandBuffer, vertexCount uint32, instanceCount uint32, firstVertex uint32, firstInstance uint32) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cvertexCount, _ := (C.uint32_t)(vertexCount), cgoAllocsUnknown
	cinstanceCount, _ := (C.uint32_t)(instanceCount), cgoAllocsUnknown
	cfirstVertex, _ := (C.uint32_t)(firstVertex), cgoAllocsUnknown
	cfirstInstance, _ := (C.uint32_t)(firstInstance), cgoAllocsUnknown
	C.callVkCmdDraw(ccommandBuffer, cvertexCount, cinstanceCount, cfirstVertex, cfirstInstance)
}

// CmdSetScissor function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdSetScissor.html
func CmdSetScissor(commandBuffer CommandBuffer, firstScissor uint32, scissorCount uint32, pScissors []Rect2D) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cfirstScissor, _ := (C.uint32_t)(firstScissor), cgoAllocsUnknown
	cscissorCount, _ := (C.uint32_t)(scissorCount), cgoAllocsUnknown
	cpScissors, _ := unpackArgSRect2D(pScissors)
	C.callVkCmdSetScissor(ccommandBuffer, cfirstScissor, cscissorCount, cpScissors)
	packSRect2D(pScissors, cpScissors)
}

// CmdSetViewport function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdSetViewport.html
func CmdSetViewport(commandBuffer CommandBuffer, firstViewport uint32, viewportCount uint32, pViewports []Viewport) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cfirstViewport, _ := (C.uint32_t)(firstViewport), cgoAllocsUnknown
	cviewportCount, _ := (C.uint32_t)(viewportCount), cgoAllocsUnknown
	cpViewports, _ := unpackArgSViewport(pViewports)
	C.callVkCmdSetViewport(ccommandBuffer, cfirstViewport, cviewportCount, cpViewports)
	packSViewport(pViewports, cpViewports)
}

// CmdBindDescriptorSets function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdBindDescriptorSets.html
func CmdBindDescriptorSets(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, layout PipelineLayout, firstSet uint32, descriptorSetCount uint32, pDescriptorSets []DescriptorSet, dynamicOffsetCount uint32, pDynamicOffsets []uint32) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cpipelineBindPoint, _ := (C.VkPipelineBindPoint)(pipelineBindPoint), cgoAllocsUnknown
	clayout, _ := *(*C.VkPipelineLayout)(unsafe.Pointer(&layout)), cgoAllocsUnknown
	cfirstSet, _ := (C.uint32_t)(firstSet), cgoAllocsUnknown
	cdescriptorSetCount, _ := (C.uint32_t)(descriptorSetCount), cgoAllocsUnknown
	cpDescriptorSets, _ := (*C.VkDescriptorSet)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pDescriptorSets)).Data)), cgoAllocsUnknown
	cdynamicOffsetCount, _ := (C.uint32_t)(dynamicOffsetCount), cgoAllocsUnknown
	cpDynamicOffsets, _ := (*C.uint32_t)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pDynamicOffsets)).Data)), cgoAllocsUnknown
	C.callVkCmdBindDescriptorSets(ccommandBuffer, cpipelineBindPoint, clayout, cfirstSet, cdescriptorSetCount, cpDescriptorSets, cdynamicOffsetCount, cpDynamicOffsets)
}

// CmdBindPipeline function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdBindPipeline.html
func CmdBindPipeline(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, pipeline Pipeline) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	cpipelineBindPoint, _ := (C.VkPipelineBindPoint)(pipelineBindPoint), cgoAllocsUnknown
	cpipeline, _ := *(*C.VkPipeline)(unsafe.Pointer(&pipeline)), cgoAllocsUnknown
	C.callVkCmdBindPipeline(ccommandBuffer, cpipelineBindPoint, cpipeline)
}

// CmdEndRenderPass function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCmdEndRenderPass.html
func CmdEndRenderPass(commandBuffer CommandBuffer) {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	C.callVkCmdEndRenderPass(ccommandBuffer)
}

// EndCommandBuffer function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkEndCommandBuffer.html
func EndCommandBuffer(commandBuffer CommandBuffer) Result {
	ccommandBuffer, _ := *(*C.VkCommandBuffer)(unsafe.Pointer(&commandBuffer)), cgoAllocsUnknown
	__ret := C.callVkEndCommandBuffer(ccommandBuffer)
	__v := (Result)(__ret)
	return __v
}

// QueueSubmit function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkQueueSubmit.html
func QueueSubmit(queue Queue, submitCount uint32, pSubmits []SubmitInfo, fence Fence) Result {
	cqueue, _ := *(*C.VkQueue)(unsafe.Pointer(&queue)), cgoAllocsUnknown
	csubmitCount, _ := (C.uint32_t)(submitCount), cgoAllocsUnknown
	cpSubmits, _ := unpackArgSSubmitInfo(pSubmits)
	cfence, _ := *(*C.VkFence)(unsafe.Pointer(&fence)), cgoAllocsUnknown
	__ret := C.callVkQueueSubmit(cqueue, csubmitCount, cpSubmits, cfence)
	packSSubmitInfo(pSubmits, cpSubmits)
	__v := (Result)(__ret)
	return __v
}

// CreateImageView function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateImageView.html
func CreateImageView(device Device, pCreateInfo *ImageViewCreateInfo, pAllocator *AllocationCallbacks, pView *ImageView) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpView, _ := (*C.VkImageView)(unsafe.Pointer(pView)), cgoAllocsUnknown
	__ret := C.callVkCreateImageView(cdevice, cpCreateInfo, cpAllocator, cpView)
	__v := (Result)(__ret)
	return __v
}

// CreateFramebuffer function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateFramebuffer.html
func CreateFramebuffer(device Device, pCreateInfo *FramebufferCreateInfo, pAllocator *AllocationCallbacks, pFramebuffer *Framebuffer) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpFramebuffer, _ := (*C.VkFramebuffer)(unsafe.Pointer(pFramebuffer)), cgoAllocsUnknown
	__ret := C.callVkCreateFramebuffer(cdevice, cpCreateInfo, cpAllocator, cpFramebuffer)
	__v := (Result)(__ret)
	return __v
}

// CreateFence function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateFence.html
func CreateFence(device Device, pCreateInfo *FenceCreateInfo, pAllocator *AllocationCallbacks, pFence *Fence) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpFence, _ := (*C.VkFence)(unsafe.Pointer(pFence)), cgoAllocsUnknown
	__ret := C.callVkCreateFence(cdevice, cpCreateInfo, cpAllocator, cpFence)
	__v := (Result)(__ret)
	return __v
}

// AllocateCommandBuffers function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkAllocateCommandBuffers.html
func AllocateCommandBuffers(device Device, pAllocateInfo *CommandBufferAllocateInfo, pCommandBuffers []CommandBuffer) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpAllocateInfo, _ := pAllocateInfo.PassRef()
	cpCommandBuffers, _ := (*C.VkCommandBuffer)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pCommandBuffers)).Data)), cgoAllocsUnknown
	__ret := C.callVkAllocateCommandBuffers(cdevice, cpAllocateInfo, cpCommandBuffers)
	__v := (Result)(__ret)
	return __v
}

// GetPhysicalDeviceSurfaceCapabilities function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkGetPhysicalDeviceSurfaceCapabilitiesKHR
func GetPhysicalDeviceSurfaceCapabilities(physicalDevice PhysicalDevice, surface Surface, pSurfaceCapabilities *SurfaceCapabilities) Result {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	csurface, _ := *(*C.VkSurfaceKHR)(unsafe.Pointer(&surface)), cgoAllocsUnknown
	cpSurfaceCapabilities, _ := pSurfaceCapabilities.PassRef()
	__ret := C.callVkGetPhysicalDeviceSurfaceCapabilitiesKHR(cphysicalDevice, csurface, cpSurfaceCapabilities)

	pSurfaceCapabilities.SupportedCompositeAlpha = CompositeAlphaFlags(cpSurfaceCapabilities.supportedCompositeAlpha)
	pSurfaceCapabilities.MinImageCount = uint32(cpSurfaceCapabilities.minImageCount)
	pSurfaceCapabilities.MaxImageCount = uint32(cpSurfaceCapabilities.maxImageCount)

	__v := (Result)(__ret)
	return __v
}

// GetPhysicalDeviceSurfaceSupport function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkGetPhysicalDeviceSurfaceSupportKHR
func GetPhysicalDeviceSurfaceSupport(physicalDevice PhysicalDevice, queueFamilyIndex uint32, surface Surface, pSupported *Bool32) Result {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cqueueFamilyIndex, _ := (C.uint32_t)(queueFamilyIndex), cgoAllocsUnknown
	csurface, _ := *(*C.VkSurfaceKHR)(unsafe.Pointer(&surface)), cgoAllocsUnknown
	cpSupported, _ := (*C.VkBool32)(unsafe.Pointer(pSupported)), cgoAllocsUnknown
	__ret := C.callVkGetPhysicalDeviceSurfaceSupportKHR(cphysicalDevice, cqueueFamilyIndex, csurface, cpSupported)
	__v := (Result)(__ret)
	return __v
}

// GetPhysicalDeviceSurfacePresentModes function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkGetPhysicalDeviceSurfacePresentModesKHR
func GetPhysicalDeviceSurfacePresentModes(physicalDevice PhysicalDevice, surface Surface, pPresentModeCount *uint32, pPresentModes []PresentMode) Result {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	csurface, _ := *(*C.VkSurfaceKHR)(unsafe.Pointer(&surface)), cgoAllocsUnknown
	cpPresentModeCount, _ := (*C.uint32_t)(unsafe.Pointer(pPresentModeCount)), cgoAllocsUnknown
	cpPresentModes, _ := (*C.VkPresentModeKHR)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pPresentModes)).Data)), cgoAllocsUnknown
	__ret := C.callVkGetPhysicalDeviceSurfacePresentModesKHR(cphysicalDevice, csurface, cpPresentModeCount, cpPresentModes)
	__v := (Result)(__ret)
	return __v
}

// CreateSwapchain function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkCreateSwapchainKHR
func CreateSwapchain(device Device, pCreateInfo *SwapchainCreateInfo, pAllocator *AllocationCallbacks, pSwapchain *Swapchain) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpSwapchain, _ := (*C.VkSwapchainKHR)(unsafe.Pointer(pSwapchain)), cgoAllocsUnknown
	__ret := C.callVkCreateSwapchainKHR(cdevice, cpCreateInfo, cpAllocator, cpSwapchain)
	__v := (Result)(__ret)
	return __v
}

// GetSwapchainImages function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkGetSwapchainImagesKHR
func GetSwapchainImages(device Device, swapchain Swapchain, pSwapchainImageCount *uint32, pSwapchainImages []Image) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cswapchain, _ := *(*C.VkSwapchainKHR)(unsafe.Pointer(&swapchain)), cgoAllocsUnknown
	cpSwapchainImageCount, _ := (*C.uint32_t)(unsafe.Pointer(pSwapchainImageCount)), cgoAllocsUnknown
	cpSwapchainImages, _ := (*C.VkImage)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pSwapchainImages)).Data)), cgoAllocsUnknown
	__ret := C.callVkGetSwapchainImagesKHR(cdevice, cswapchain, cpSwapchainImageCount, cpSwapchainImages)
	__v := (Result)(__ret)
	return __v
}

// GetPhysicalDeviceSurfaceFormats function as declared in https://www.khronos.org/registry/vulkan/specs/1.0-wsi_extensions/xhtml/vkspec.html#VkGetPhysicalDeviceSurfaceFormatsKHR
func GetPhysicalDeviceSurfaceFormats(physicalDevice PhysicalDevice, surface Surface, pSurfaceFormatCount *uint32, pSurfaceFormats []SurfaceFormat) Result {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	csurface, _ := *(*C.VkSurfaceKHR)(unsafe.Pointer(&surface)), cgoAllocsUnknown
	cpSurfaceFormatCount, _ := (*C.uint32_t)(unsafe.Pointer(pSurfaceFormatCount)), cgoAllocsUnknown
	cpSurfaceFormats, _ := unpackArgSSurfaceFormat(pSurfaceFormats)
	__ret := C.callVkGetPhysicalDeviceSurfaceFormatsKHR(cphysicalDevice, csurface, cpSurfaceFormatCount, cpSurfaceFormats)
	packSSurfaceFormat(pSurfaceFormats, cpSurfaceFormats)
	__v := (Result)(__ret)
	for i := range pSurfaceFormats {
		pSurfaceFormats[i].Format = Format(pSurfaceFormats[i].refedaf82ca.format)
		pSurfaceFormats[i].ColorSpace = ColorSpace(pSurfaceFormats[i].refedaf82ca.colorSpace)
	}
	return __v
}

// CreateRenderPass function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateRenderPass.html
func CreateRenderPass(device Device, pCreateInfo *RenderPassCreateInfo, pAllocator *AllocationCallbacks, pRenderPass *RenderPass) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpRenderPass, _ := (*C.VkRenderPass)(unsafe.Pointer(pRenderPass)), cgoAllocsUnknown
	__ret := C.callVkCreateRenderPass(cdevice, cpCreateInfo, cpAllocator, cpRenderPass)
	__v := (Result)(__ret)
	return __v
}

// CreateCommandPool function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateCommandPool.html
func CreateCommandPool(device Device, pCreateInfo *CommandPoolCreateInfo, pAllocator *AllocationCallbacks, pCommandPool *CommandPool) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpCommandPool, _ := (*C.VkCommandPool)(unsafe.Pointer(pCommandPool)), cgoAllocsUnknown
	__ret := C.callVkCreateCommandPool(cdevice, cpCreateInfo, cpAllocator, cpCommandPool)
	__v := (Result)(__ret)
	return __v
}

// CreateSemaphore function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateSemaphore.html
func CreateSemaphore(device Device, pCreateInfo *SemaphoreCreateInfo, pAllocator *AllocationCallbacks, pSemaphore *Semaphore) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpSemaphore, _ := (*C.VkSemaphore)(unsafe.Pointer(pSemaphore)), cgoAllocsUnknown
	__ret := C.callVkCreateSemaphore(cdevice, cpCreateInfo, cpAllocator, cpSemaphore)
	__v := (Result)(__ret)
	return __v
}

// CreateDescriptorSetLayout function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateDescriptorSetLayout.html
func CreateDescriptorSetLayout(device Device, pCreateInfo *DescriptorSetLayoutCreateInfo, pAllocator *AllocationCallbacks, pSetLayout *DescriptorSetLayout) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpSetLayout, _ := (*C.VkDescriptorSetLayout)(unsafe.Pointer(pSetLayout)), cgoAllocsUnknown
	__ret := C.callVkCreateDescriptorSetLayout(cdevice, cpCreateInfo, cpAllocator, cpSetLayout)
	__v := (Result)(__ret)
	return __v
}

// CreatePipelineLayout function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreatePipelineLayout.html
func CreatePipelineLayout(device Device, pCreateInfo *PipelineLayoutCreateInfo, pAllocator *AllocationCallbacks, pPipelineLayout *PipelineLayout) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpPipelineLayout, _ := (*C.VkPipelineLayout)(unsafe.Pointer(pPipelineLayout)), cgoAllocsUnknown
	__ret := C.callVkCreatePipelineLayout(cdevice, cpCreateInfo, cpAllocator, cpPipelineLayout)
	__v := (Result)(__ret)
	return __v
}

// CreateShaderModule function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateShaderModule.html
func CreateShaderModule(device Device, pCreateInfo *ShaderModuleCreateInfo, pAllocator *AllocationCallbacks, pShaderModule *ShaderModule) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpShaderModule, _ := (*C.VkShaderModule)(unsafe.Pointer(pShaderModule)), cgoAllocsUnknown
	__ret := C.callVkCreateShaderModule(cdevice, cpCreateInfo, cpAllocator, cpShaderModule)
	__v := (Result)(__ret)
	return __v
}

// CreateGraphicsPipelines function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateGraphicsPipelines.html
func CreateGraphicsPipelines(device Device, pipelineCache PipelineCache, createInfoCount uint32, pCreateInfos []GraphicsPipelineCreateInfo, pAllocator *AllocationCallbacks, pPipelines []Pipeline) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpipelineCache, _ := *(*C.VkPipelineCache)(unsafe.Pointer(&pipelineCache)), cgoAllocsUnknown
	ccreateInfoCount, _ := (C.uint32_t)(createInfoCount), cgoAllocsUnknown
	cpCreateInfos, _ := unpackArgSGraphicsPipelineCreateInfo(pCreateInfos)
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpPipelines, _ := (*C.VkPipeline)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pPipelines)).Data)), cgoAllocsUnknown
	__ret := C.callVkCreateGraphicsPipelines(cdevice, cpipelineCache, ccreateInfoCount, cpCreateInfos, cpAllocator, cpPipelines)
	packSGraphicsPipelineCreateInfo(pCreateInfos, cpCreateInfos)
	__v := (Result)(__ret)
	return __v
}

// CreateBuffer function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateBuffer.html
func CreateBuffer(device Device, pCreateInfo *BufferCreateInfo, pAllocator *AllocationCallbacks, pBuffer *Buffer) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpBuffer, _ := (*C.VkBuffer)(unsafe.Pointer(pBuffer)), cgoAllocsUnknown
	__ret := C.callVkCreateBuffer(cdevice, cpCreateInfo, cpAllocator, cpBuffer)
	__v := (Result)(__ret)
	return __v
}

// GetBufferMemoryRequirements function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetBufferMemoryRequirements.html
func GetBufferMemoryRequirements(device Device, buffer Buffer, pMemoryRequirements *MemoryRequirements) {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cbuffer, _ := *(*C.VkBuffer)(unsafe.Pointer(&buffer)), cgoAllocsUnknown
	cpMemoryRequirements, _ := pMemoryRequirements.PassRef()
	C.callVkGetBufferMemoryRequirements(cdevice, cbuffer, cpMemoryRequirements)

	pMemoryRequirements.MemoryTypeBits = uint32(cpMemoryRequirements.memoryTypeBits)
}

// AllocateMemory function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkAllocateMemory.html
func AllocateMemory(device Device, pAllocateInfo *MemoryAllocateInfo, pAllocator *AllocationCallbacks, pMemory *DeviceMemory) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpAllocateInfo, _ := pAllocateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpMemory, _ := (*C.VkDeviceMemory)(unsafe.Pointer(pMemory)), cgoAllocsUnknown
	__ret := C.callVkAllocateMemory(cdevice, cpAllocateInfo, cpAllocator, cpMemory)
	__v := (Result)(__ret)
	return __v
}

// MapMemory function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkMapMemory.html
func MapMemory(device Device, memory DeviceMemory, offset DeviceSize, size DeviceSize, flags MemoryMapFlags, ppData *unsafe.Pointer) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cmemory, _ := *(*C.VkDeviceMemory)(unsafe.Pointer(&memory)), cgoAllocsUnknown
	coffset, _ := (C.VkDeviceSize)(offset), cgoAllocsUnknown
	csize, _ := (C.VkDeviceSize)(size), cgoAllocsUnknown
	cflags, _ := (C.VkMemoryMapFlags)(flags), cgoAllocsUnknown
	cppData, _ := ppData, cgoAllocsUnknown
	__ret := C.callVkMapMemory(cdevice, cmemory, coffset, csize, cflags, cppData)
	__v := (Result)(__ret)
	return __v
}

// BindBufferMemory function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkBindBufferMemory.html
func BindBufferMemory(device Device, buffer Buffer, memory DeviceMemory, memoryOffset DeviceSize) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cbuffer, _ := *(*C.VkBuffer)(unsafe.Pointer(&buffer)), cgoAllocsUnknown
	cmemory, _ := *(*C.VkDeviceMemory)(unsafe.Pointer(&memory)), cgoAllocsUnknown
	cmemoryOffset, _ := (C.VkDeviceSize)(memoryOffset), cgoAllocsUnknown
	__ret := C.callVkBindBufferMemory(cdevice, cbuffer, cmemory, cmemoryOffset)
	__v := (Result)(__ret)
	return __v
}

// CreateDescriptorPool function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateDescriptorPool.html
func CreateDescriptorPool(device Device, pCreateInfo *DescriptorPoolCreateInfo, pAllocator *AllocationCallbacks, pDescriptorPool *DescriptorPool) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpDescriptorPool, _ := (*C.VkDescriptorPool)(unsafe.Pointer(pDescriptorPool)), cgoAllocsUnknown
	__ret := C.callVkCreateDescriptorPool(cdevice, cpCreateInfo, cpAllocator, cpDescriptorPool)
	__v := (Result)(__ret)
	return __v
}

// AllocateDescriptorSets function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkAllocateDescriptorSets.html
func AllocateDescriptorSets(device Device, pAllocateInfo *DescriptorSetAllocateInfo, pDescriptorSets *DescriptorSet) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpAllocateInfo, _ := pAllocateInfo.PassRef()
	cpDescriptorSets, _ := (*C.VkDescriptorSet)(unsafe.Pointer(pDescriptorSets)), cgoAllocsUnknown
	__ret := C.callVkAllocateDescriptorSets(cdevice, cpAllocateInfo, cpDescriptorSets)
	__v := (Result)(__ret)
	return __v
}

// UpdateDescriptorSets function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkUpdateDescriptorSets.html
func UpdateDescriptorSets(device Device, descriptorWriteCount uint32, pDescriptorWrites []WriteDescriptorSet, descriptorCopyCount uint32, pDescriptorCopies []CopyDescriptorSet) {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cdescriptorWriteCount, _ := (C.uint32_t)(descriptorWriteCount), cgoAllocsUnknown
	cpDescriptorWrites, _ := unpackArgSWriteDescriptorSet(pDescriptorWrites)
	cdescriptorCopyCount, _ := (C.uint32_t)(descriptorCopyCount), cgoAllocsUnknown
	cpDescriptorCopies, _ := unpackArgSCopyDescriptorSet(pDescriptorCopies)
	C.callVkUpdateDescriptorSets(cdevice, cdescriptorWriteCount, cpDescriptorWrites, cdescriptorCopyCount, cpDescriptorCopies)
	packSCopyDescriptorSet(pDescriptorCopies, cpDescriptorCopies)
	packSWriteDescriptorSet(pDescriptorWrites, cpDescriptorWrites)
}

// CreatePipelineCache function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreatePipelineCache.html
func CreatePipelineCache(device Device, pCreateInfo *PipelineCacheCreateInfo, pAllocator *AllocationCallbacks, pPipelineCache *PipelineCache) Result {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpPipelineCache, _ := (*C.VkPipelineCache)(unsafe.Pointer(pPipelineCache)), cgoAllocsUnknown
	__ret := C.callVkCreatePipelineCache(cdevice, cpCreateInfo, cpAllocator, cpPipelineCache)
	__v := (Result)(__ret)
	return __v
}

// CreateInstance function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateInstance.html
func CreateInstance(pCreateInfo *InstanceCreateInfo, pAllocator *AllocationCallbacks, pInstance *Instance) Result {
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpInstance, _ := (*C.VkInstance)(unsafe.Pointer(pInstance)), cgoAllocsUnknown
	__ret := C.callVkCreateInstance(cpCreateInfo, cpAllocator, cpInstance)
	__v := (Result)(__ret)
	return __v
}

// EnumeratePhysicalDevices function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkEnumeratePhysicalDevices.html
func EnumeratePhysicalDevices(instance Instance, pPhysicalDeviceCount *uint32, pPhysicalDevices []PhysicalDevice) Result {
	cinstance, _ := *(*C.VkInstance)(unsafe.Pointer(&instance)), cgoAllocsUnknown
	cpPhysicalDeviceCount, _ := (*C.uint32_t)(unsafe.Pointer(pPhysicalDeviceCount)), cgoAllocsUnknown
	cpPhysicalDevices, _ := (*C.VkPhysicalDevice)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&pPhysicalDevices)).Data)), cgoAllocsUnknown
	__ret := C.callVkEnumeratePhysicalDevices(cinstance, cpPhysicalDeviceCount, cpPhysicalDevices)
	__v := (Result)(__ret)
	return __v
}

// CreateDevice function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkCreateDevice.html
func CreateDevice(physicalDevice PhysicalDevice, pCreateInfo *DeviceCreateInfo, pAllocator *AllocationCallbacks, pDevice *Device) Result {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cpCreateInfo, _ := pCreateInfo.PassRef()
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpDevice, _ := (*C.VkDevice)(unsafe.Pointer(pDevice)), cgoAllocsUnknown
	__ret := C.callVkCreateDevice(cphysicalDevice, cpCreateInfo, cpAllocator, cpDevice)
	__v := (Result)(__ret)

	*pDevice = *(*Device)(cpDevice)
	print(":")
	println(cpDevice)
	println(*pDevice)
	return __v
}

// GetDeviceQueue function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetDeviceQueue.html
func GetDeviceQueue(device Device, queueFamilyIndex uint32, queueIndex uint32, pQueue *Queue) {
	cdevice, _ := *(*C.VkDevice)(unsafe.Pointer(&device)), cgoAllocsUnknown
	cqueueFamilyIndex, _ := (C.uint32_t)(queueFamilyIndex), cgoAllocsUnknown
	cqueueIndex, _ := (C.uint32_t)(queueIndex), cgoAllocsUnknown
	cpQueue, _ := (*C.VkQueue)(unsafe.Pointer(pQueue)), cgoAllocsUnknown
	C.callVkGetDeviceQueue(cdevice, cqueueFamilyIndex, cqueueIndex, cpQueue)
	*pQueue = *(*Queue)(cpQueue)
}

// GetPhysicalDeviceQueueFamilyProperties function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetPhysicalDeviceQueueFamilyProperties.html
func GetPhysicalDeviceQueueFamilyProperties(physicalDevice PhysicalDevice, pQueueFamilyPropertyCount *uint32, pQueueFamilyProperties []QueueFamilyProperties) {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cpQueueFamilyPropertyCount, _ := (*C.uint32_t)(unsafe.Pointer(pQueueFamilyPropertyCount)), cgoAllocsUnknown
	cpQueueFamilyProperties, _ := unpackArgSQueueFamilyProperties(pQueueFamilyProperties)
	C.callVkGetPhysicalDeviceQueueFamilyProperties(cphysicalDevice, cpQueueFamilyPropertyCount, cpQueueFamilyProperties)
	packSQueueFamilyProperties(pQueueFamilyProperties, cpQueueFamilyProperties)
	for i := range pQueueFamilyProperties {
		pQueueFamilyProperties[i].Deref()
	}
}

// GetPhysicalDeviceMemoryProperties function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetPhysicalDeviceMemoryProperties.html
func GetPhysicalDeviceMemoryProperties(physicalDevice PhysicalDevice, pMemoryProperties *PhysicalDeviceMemoryProperties) {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cpMemoryProperties, _ := pMemoryProperties.PassRef()
	C.callVkGetPhysicalDeviceMemoryProperties(cphysicalDevice, cpMemoryProperties)

	pMemoryProperties.Deref()

	for i := range pMemoryProperties.MemoryTypes {
		pMemoryProperties.MemoryTypes[i].Deref()
	}
}

// GetPhysicalDeviceProperties function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetPhysicalDeviceProperties.html
func GetPhysicalDeviceProperties(physicalDevice PhysicalDevice, pProperties *PhysicalDeviceProperties) {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cpProperties, _ := pProperties.PassRef()
	C.callVkGetPhysicalDeviceProperties(cphysicalDevice, cpProperties)

	pProperties.Deref()
}

// GetPhysicalDeviceFeatures2 function as declared in https://www.khronos.org/registry/vulkan/specs/1.0/man/html/vkGetPhysicalDeviceFeatures2.html
func GetPhysicalDeviceFeatures2(physicalDevice PhysicalDevice, pFeatures *PhysicalDeviceFeatures2) {
	cphysicalDevice, _ := *(*C.VkPhysicalDevice)(unsafe.Pointer(&physicalDevice)), cgoAllocsUnknown
	cpFeatures, _ := pFeatures.PassRef()
	C.callVkGetPhysicalDeviceFeatures2(cphysicalDevice, cpFeatures)
}

// Wayland-related

func CreateWaylandSurface(instance Instance, info *WaylandSurfaceCreateInfo, pAllocator *AllocationCallbacks, pSurface *Surface) {
	cpAllocator, _ := (*C.VkAllocationCallbacks)(unsafe.Pointer(pAllocator)), cgoAllocsUnknown
	cpSurface, _ := (*C.VkSurfaceKHR)(unsafe.Pointer(pSurface)), cgoAllocsUnknown

	C.callCreateWaylandSurfaceKHR(instance, (*C.VkWaylandSurfaceCreateInfoKHR)(unsafe.Pointer(info)), cpAllocator, cpSurface)
}

func GetPhysicalDeviceWaylandPresentationSupport(instance Instance, physicalDevice PhysicalDevice, queueFamilyIndex uint32, display uintptr) bool {
	return 0 != C.callGetPhysicalDeviceWaylandPresentationSupportKHR(instance, physicalDevice, C.uint(queueFamilyIndex), (*C.struct_wl_display)(unsafe.Pointer(display)))
}
