<template>
    <div class="json-details-layout">
        <pre>{{details}}</pre>
    </div>
</template>

<script>
    export default {
        props: {
            details: Object
        }
    }
</script>

<style lang="scss" scoped>
    .json-details-layout {
        margin: 20px 40px;
        padding: 20px;
        background: #455070;
        color: #bfc6e0;
        border-radius: 2px;
        font-size: 12px;
        line-height: 20px;
        height: calc(100% - 40px);
        overflow: auto;
    }
</style>
