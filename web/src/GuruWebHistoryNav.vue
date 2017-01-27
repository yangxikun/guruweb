<template>
	<div>
		<h3 class="left">{{ currentFile }}</h3>
		<div class="cursor-pointer" v-show="historyStack.length > 1" @click="historyBack()">
			<h3 class="right">back</h3>
			<span class="right glyphicon glyphicon-menu-left"></span>
		</div>
	</div>
</template>
<style scoped>
	.left {
		float: left;
	}
	.right {
		float: right;
	}
	.cursor-pointer {
		cursor: pointer;
	}
	span {
	    font-size: 24px;
	    margin-top: 20px;
	    margin-bottom: 10px;
	}
</style>
<script>
    export default {
        name: 'history',
        data () {
            return {
                currentFile: '',
                historyStack: []
            };
        },
        mounted() {
        	window.Bus.$on('show-file', (fileInfo) => {
        		this.currentFile = fileInfo.fileName;
        		this.historyStack.push(fileInfo)
        	});
        },
        methods: {
            historyBack() {
            	this.historyStack.pop(); // pop current
            	window.Bus.$emit('show-file', this.historyStack.pop());
            }
        }
    }
</script>