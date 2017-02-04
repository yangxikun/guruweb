<template>
    <div class="menu" v-bind:style="{ top: top + 'px', left: left + 'px' }" v-show="ok">
        <ul class="list-group">
            <li class="list-group-item" v-for="item in menu" v-on:click="query(item)">{{ item }}</li>
        </ul>
    </div>
</template>
<style lang="less">
    .menu {
        position: fixed;
        .list-group-item {
            cursor: pointer;
        }
    }
</style>
<script>
    export default {
        name: 'menu',
        data() {
            return{
                ok: false,
                menu: new Array(),
                pos: '',
                top: 0,
                left: 0
            }
        },
        mounted() {
            window.Bus.$on('toggle-menu', (pos, top, left) => {
                console.log('show menu');
                this.ok = !this.ok;
                if (this.ok) {
                    this.$http.get('/query?mode=what&pos=' + encodeURIComponent(pos), {
                        headers: {Accept: 'application/json'}
                    }).then((response) => {
                        this.menu = response.data.modes;
                        this.top = top;
                        this.left = left;
                        this.pos = pos;
                    }, (response) => {
                        window.Bus.$emit('show-alert', response.data);
                    });
                }
            });
            window.Bus.$on('change-body', (name) => {
                if (this.ok) {
                this.ok = false;
            }
        });
        },
        methods: {
            query(mode) {
                this.ok = false;
                console.log('query');
                console.log(mode);
                window.Bus.$emit('show-query-out', mode, this.pos);
            }
        }
    }
</script>
