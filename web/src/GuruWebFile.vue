<template>
    <div class="col-lg-8 code" id="code" :style="{height: height + 'px'}">
        <div class="row">
            <div class="col-lg-1 nums">
                <template v-for="num in nums">
                    <span :id="'L' + num">{{ num }}</span>
                    <br>
                </template>
            </div>
            <div class="col-lg-11 lines" @mouseup="what"></div>
        </div>
    </div>
</template>
<style>
    .code {
        background: #e9e9e9;
        border-radius: 5px;
        font: 14px Menlo, monospace;
        padding: 10px;
        overflow: auto;
    }
    .nums {
        color: #999;
        cursor: default;
        padding-right: 10px;
        text-align: right;
        width: 48px;
    }
    .lines {
        white-space: pre;
    }
    .comment {
        color: #060;
    }
    .selection {
        background: #ff9632;
    }
    .row {
        margin: 0px 0px;
    }
</style>
<script>
    import Vue from 'vue';
    export default {
        name: 'file',
        data() {
            return {
                nums: 0,
                currentFile: '',
                height: 0
            };
        },
        mounted() {
            const topOuterHeight = $('#top').outerHeight(true) + 61;
            $(window).resize(() => {
                this.height = $(window).height() - topOuterHeight;
            });
            window.Bus.$on('show-file', (fileInfo) => {
                console.log('show file', fileInfo)
                let params = 'path=' + encodeURIComponent(fileInfo.fileName);
                if (fileInfo.sel !== null) {
                    params = params + '&s=' + fileInfo.sel.fromLine + '.' + fileInfo.sel.fromCol + '-'
                        + fileInfo.sel.toLine + '.' + fileInfo.sel.toCol;
                }
                this.currentFile = fileInfo.fileName;
                this.$http.get('/file?' + params, {}).then((response) => {
                    console.log('show file http')
                    this.nums = (response.data.match(/\n/g) || []).length;
                    $('.lines').html(response.data);
                    this.height = $(window).height() - topOuterHeight;
                    Vue.nextTick(function () {
                        if (fileInfo.sel != null) {
                            const l = $('#L' + fileInfo.sel.fromLine);
                            if (l.offset().top < topOuterHeight) {
                                $('#code').scrollTop(
                                    $('#code').scrollTop() - (topOuterHeight - l.offset().top + 3*l.height())
                                );
                            } else {
                                $('#code').scrollTop(
                                    $('#code').scrollTop() + (l.offset().top - topOuterHeight - 3*l.height())
                                );
                            }
                        }
                    });
                }, (response) => {
                    window.Bus.$emit('show-alert', response.data);
                });
          });
        },
        methods: {
            what(event) {
                console.log('mouse up');
                let lines = $('.lines');
                let range = window.getSelection().getRangeAt(0);
                this._insertSelectionMarks(range);
                let sel = this._selectionMarkOffsets($(lines));
                lines.find('.mark').detach();
                let b = this._getByteOffsets(lines.text(), sel);
                let pos = this._pos(this.currentFile, b.startOffset, b.endOffset);
                window.Bus.$emit('toggle-menu', pos, event.clientY, event.clientX);
            },
            _pos(file, start, end) {
                let p = file + ':#' + start;
                if (start != end) {
                    p += ',#' + end;
                }
                return p;
            },
            // From http://stackoverflow.com/a/12206089
            _getUTF8Length(s, start, end) {
                let len = 0;
                for (let i = start; i < end; i++) {
                    let code = s.charCodeAt(i);
                    if (code <= 0x7f) {
                        len += 1;
                    } else if (code <= 0x7ff) {
                        len += 2;
                    } else if (code >= 0xd800 && code <= 0xdfff) {
                        // Surrogate pair: These take 4 bytes in UTF-8 and 2 chars in UCS-2
                        // (Assume next char is the other [valid] half and just skip it)
                        len += 4;
                        i++;
                    } else if (code < 0xffff) {
                        len += 3;
                    } else {
                        len += 4;
                    }
                }
                return len;
            },
            _getByteOffsets(s, range) {
                const a = this._getUTF8Length(s, 0, range.startOffset);
                const b = this._getUTF8Length(s, range.startOffset, range.endOffset);
                return {startOffset: a, endOffset: a + b};
            },
            _selectionMarkOffsets(lines) {
                const marked = lines.text();
                return {
                    startOffset: marked.indexOf('\0'),
                    endOffset: marked.lastIndexOf('\0') - 1
                };
            },
            _insertSelectionMarks(range) {
                let s = range.cloneRange();
                let e = range.cloneRange();
                s.collapse(true);
                e.collapse(false);
                s.insertNode(this._makeSelectionMark('start'));
                e.insertNode(this._makeSelectionMark('end'));
            },
            _makeSelectionMark(type) {
              return $('<span class="mark">').addClass(type).text('\0')[0];
            }
        }
    }
</script>
