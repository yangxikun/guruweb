<template>
    <nav class="navbar navbar-default" id="top">
        <div class="container-fluid">
            <div class="row">
                <div class="col-lg-4"></div>
                <div class="col-lg-4">
                    <a class="navbar-brand" href="#">GuruWeb</a>
                    <ul class="nav navbar-nav">
                        <li :class="{active: configActive}">
                            <a href="#" @click="showConfig()">Config</a>
                        </li>
                    </ul>

                    <form class="navbar-form navbar-right">
                        <div class="form-group">
                            <input v-model="searchString" type="text" class="form-control" placeholder="Search"
                                @focus="recommendSearch()" @input="recommendSearch()">
                        </div>
                    </form>
                </div>
                <div class="col-lg-4"></div>
            </div>
        </div>
    </nav>
</template>
<style>
</style>
<script>
    export default {
        name: 'head',
        data () {
            return {
                searchString: '',
                configActive: false
            };
        },
        methods: {
            showConfig() {
                this.configActive = true;
                window.Bus.$emit('change-body', 'guru-web-config')
            },
            recommendSearch() {
                const val = this.searchString;
                if (val.length >= 2) {
                    this.configActive = false;
                    console.log('guru-web-head', val)
                    window.Bus.$emit('search-change', val);
                    window.Bus.$emit('change-body', 'guru-web-search');
                }
            }
        }
    }
</script>
