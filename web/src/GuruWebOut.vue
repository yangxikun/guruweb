<template>
    <div class="col-lg-4 query-out" :style="{height: height + 'px'}">
        <template v-for="link in links">
            <a v-if="typeof link === 'object'" :title="link.tooltip" @click="jump(link.file, link.sel)">{{ link.text }}</a><span v-if="typeof link === 'string'">{{ link }}</span>
        </template>
    </div>
</template>
<style lang="less">
    .query-out {
        padding: 10px;
        white-space: pre;
        font: 14px Menlo, monospace;
        overflow: auto;
        a {
            cursor: pointer;
            text-decoration: none;
            display: block;
        }
    }
</style>
<script>
    export default {
        name: 'out',
        data() {
            return {
                links: ['Select or click within the source code to consult the guru.'],
                height: 0,
            };
        },
        mounted() {
            const topOuterHeight = $('#top').outerHeight(true) + 61;
            $(window).resize(() => {
                this.height = $(window).height() - topOuterHeight;
            });
            window.Bus.$on('show-query-out', (mode, pos) => {
                console.log('show-query-out');
                this.links = ["loading......"];
                this.$http.get('/query?mode=' + mode + '&pos=' + encodeURIComponent(pos), {}).then((response) => {
                    this.height = $(window).height() - topOuterHeight;
                    this.links = this._parseLinks(response.data);
                }, (response) => {
                    this.links = ['Select or click within the source code to consult the guru.'];
                    window.Bus.$emit('show-alert', response.data);
                });
            });
            window.Bus.$on('change-body', (name) => {
                if (name == 'guru-web-file') {
                    this.height = $(window).height() - topOuterHeight;
                    this.links = ['Select or click within the source code to consult the guru.'];
                }
            });
        },
        methods: {
            jump(file, sel) {
                console.log(file, sel);
                const fileInfo = {fileName: file, sel: sel};
                window.Bus.$emit('show-file', fileInfo);
            },
            _parseLinks(text) {
                // file:line.col-line.col:
                const rangeAddress = /(.*):([0-9]+)\.([0-9]+)-([0-9]+)\.([0-9]+): (.*)/;
                // file:line:col:
                const singleAddress = /(.*):([0-9]+):([0-9]+): (.*)/;
                // -:
                const noAddress = /-: (.*)/;
                const arrow = '▶ ';
                const lines = text.split('\n');
                const n = lines.length;
                let links = new Array();
                let match;
                for (let i = 0; i < n; i++) {
                    const line = lines[i];
                    if (match = rangeAddress.exec(line)) {
                        const file = match[1];
                        const sel = {
                            fromLine: parseInt(match[2], 10),
                            fromCol: parseInt(match[3], 10),
                            toLine: parseInt(match[4], 10),
                            toCol: parseInt(match[5], 10)
                        };
                        const rest = match[6];
                        links.push({
                            file: file,
                            sel: sel,
                            text: arrow + rest,
                            tooltip: line
                        });
                        continue;
                    }
                    if (match = singleAddress.exec(line)) {
                        const file = match[1];
                        const sel = {
                            fromLine: parseInt(match[2], 10),
                            fromCol: parseInt(match[3], 10),
                            toLine: parseInt(match[2], 10),
                            toCol: parseInt(match[3], 10)
                        };
                        const rest = match[4];
                        links.push({
                            file: file,
                            sel: sel,
                            text: arrow + rest,
                            tooltip: line
                        });
                        continue;
                    }
                    if (match = noAddress.exec(line)) {
                        const rest = match[1];
                        links.push('  ' + rest + '\n');
                        continue;
                    }
                    links.push(line + '\n');
                }
                return links;
            }
        }
    }
</script>
