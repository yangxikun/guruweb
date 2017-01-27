<template>
	<div class="alert alert-warning" v-show="ok">
		 <span class="close" @click="close">&times;</span>
		{{ content }}
	</div>
</template>
<style>
	.alert {
		position: fixed;
		left: 25%;
		top: 25%;
		width: 50%;
		max-height: 50%;
		overflow: auto;
	}
</style>
<script>
	export default {
		name: 'alert',
		data() {
			return {
				ok: false,
				content: ''
			}
		},
		mounted() {
			window.Bus.$on('show-alert', (content) => {
				this.ok = true;
				this.content = content;
			});
		},
		methods: {
			close() {
				this.ok = false;
			}
		}
	}
</script>