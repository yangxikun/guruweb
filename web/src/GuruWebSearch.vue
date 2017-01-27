<template v-show="ok">
    <div>
        <div class="clearfix">
            <h3>Files <span class="badge">{{ fileItems.length }}</span></h3>
            <a href="#" v-for="item in fileItems" @click="showFile(item)">{{ item }}</a>
        </div>
        <div class="clearfix">
            <h3>Struct Or Interface <span class="badge">{{ Object.keys(identItems).length }}</span></h3>
            <a href="#" v-for="(item, key) in identItems" @click="showFile(item)">{{ key }}</a>
        </div>
    <div>
</template>
<style scoped>
    a {
        display: block;
        float: left;
        padding: 5px 20px;
    }
    .clearfix {
        overflow: auto;
    }
</style>
<script>
    export default {
        name: 'search',
        data () {
            return{
                fileItems: [],
                identItems: {}
            }
        },
        mounted: function() {
            window.Bus.$on('search-change', (val) => {
                console.log('guru-web-search', val);
                this.$http.get('/recommend-search?query=' + val, {}).then((response) => {
                    console.log('guru-web-search', response);
                    const jsonResp = response.data;
                    this.fileItems = jsonResp.files;
                    this.identItems = jsonResp.idents;
                }, (response) => {
                    window.Bus.$emit('show-alert', response.data);
                });
          });
        },
        methods: {
            showFile(item) {
                let fileInfo = {fileName: '', sel: null};
                if (typeof item == 'string') {
                    fileInfo.fileName = item;
                } else if (typeof item == 'object') {
                    fileInfo.fileName = item.Filename;
                    fileInfo.sel = {
                        fromLine: item.Line,
                        fromCol: item.Column,
                        toLine: item.Line,
                        toCol: item.Column
                    };
                }
                console.log('guru-web-search ', item);
                window.Bus.$emit('show-file', fileInfo);
                window.Bus.$emit('change-body', 'guru-web-file');
            }
        }
    }
</script>
