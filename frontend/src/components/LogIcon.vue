<template>
  <svg 
    class="log-icon" 
    viewBox="0 0 1024 1024" 
    version="1.1" 
    xmlns="http://www.w3.org/2000/svg"
    :style="{ width: size + 'px', height: size + 'px' }"
  >
    <path d="M928 256h-32V224h-96v8.64L788.8 224H800V128h-32V96h-96v32h-32v35.52l-13.76-1.6A88.32 88.32 0 0 0 608 160h-101.44a96 96 0 0 0-181.12 0H224a96 96 0 0 0-96 96v576a96 96 0 0 0 96 96h608a96 96 0 0 0 37.44-7.68 80.96 80.96 0 0 0 16.32-8.96 94.08 94.08 0 0 0 25.92-25.6A96 96 0 0 0 928 832v-32a96 96 0 0 0-16.32-53.76 94.08 94.08 0 0 0-25.92-25.6 80.96 80.96 0 0 0-16.32-8.96 89.92 89.92 0 0 0-15.68-5.12A335.68 335.68 0 0 0 896 352h32z" fill="#FFFFFF"></path>
    <path d="M592 496m-304 0a304 304 0 1 0 608 0 304 304 0 1 0-608 0Z" fill="#E9EAEB"></path>
    <path d="M768 160h-32V128h-32v32h-32v32h32v32h32V192h32V160zM864 256h-32v32h-32v32h32v32h32v-32h32V288h-32V256z" fill="#BCC0C4"></path>
    <path d="M160 192m64 0l384 0q64 0 64 64l0 576q0 64-64 64l-384 0q-64 0-64-64l0-576q0-64 64-64Z" fill="#FFFFFF"></path>
    <path d="M448 352h-128v32h128v-32zM288 352H224v32h64v-32zM608 448h-288v32h288v-32zM288 448H224v32h64v-32zM384 544h-64v32h64v-32zM288 544H224v32h64v-32z" :fill="color"></path>
    <path d="M608 192h-32v32h32a32 32 0 0 1 32 32v192h32V256a64 64 0 0 0-64-64zM192 832V256a32 32 0 0 1 32-32h32V192H224a64 64 0 0 0-64 64v576a64 64 0 0 0 64 64h192v-32H224a32 32 0 0 1-32-32z" :fill="color"></path>
    <path d="M480 192a64 64 0 0 0-128 0H288v96h256V192z" fill="#A3D4FF"></path>
    <path d="M416 160a32 32 0 0 1 32 32v32h64v32h-192V224h64V192a32 32 0 0 1 32-32m0-32a64 64 0 0 0-64 64H288v96h256V192h-64a64 64 0 0 0-64-64z" :fill="color"></path>
    <path d="M480 736m64 0l288 0q64 0 64 64l0 32q0 64-64 64l-288 0q-64 0-64-64l0-32q0-64 64-64Z" fill="#A3D4FF"></path>
    <path d="M736 576m-64 0a64 64 0 1 0 128 0 64 64 0 1 0-128 0Z" fill="#A3D4FF"></path>
    <path d="M832 736h-128l32-96a64 64 0 1 0-29.44-7.36L672 736h-128a64 64 0 0 0-64 64v32a64 64 0 0 0 64 64h288a64 64 0 0 0 64-64v-32a64 64 0 0 0-64-64z m-96-192a32 32 0 1 1-32 32 32 32 0 0 1 32-32z m128 288a32 32 0 0 1-32 32h-288a32 32 0 0 1-32-32v-32a32 32 0 0 1 32-32h288a32 32 0 0 1 32 32zM288 704a32 32 0 1 1-32 32 32 32 0 0 1 32-32m0-32a64 64 0 1 0 64 64 64 64 0 0 0-64-64z" :fill="color"></path>
  </svg>
</template>

<script setup lang="ts">
interface Props {
  size?: number
  color?: string
}

withDefaults(defineProps<Props>(), {
  size: 16,
  color: '#2A5082'
})
</script>

<style scoped>
.log-icon {
  display: inline-block;
  vertical-align: middle;
}
</style>
