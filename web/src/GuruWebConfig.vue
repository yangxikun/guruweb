<template>
	<div class="form">
		<div>
			<h3>Loaded Packages</h3>
            <div class="loadedPkgs">
                <p v-for="item in loadedPkgs">{{ item }}</p>
            </div>
		</div>
		<div>
			<h3>Input Scope</h3>
			<p v-for="(item, index) in scopes">{{ item }}<span class="close" @click="removeItem(index)">&times;</span></p>
		</div>
		<div class="input-group">
			<input type="text" class="form-control" placeholder="add scope"  v-model="newScope"
                @keyup.enter="addScope()" @keydown.up.prevent="recommendPkgsSelect(-1)"
                @keydown.down.prevent="recommendPkgsSelect(1)">
            <ul class="list-group" id="recommend-pkgs">
                <li class="list-group-item" v-for="pkg in recommendPkgs" :class="{ active: pkg.select }">{{ pkg.item }}</li>
            </ul>
		</div>
		<div class="input-group">
			<h3>Verbose</h3>
			<input type="checkbox" class="input-checkbox" v-model="verbose">
		</div>
		<button type="submit" class="btn btn-default" :class="{'btn-warning': state == 'loading'}" @click="reload()">{{ state }}</button>
	</div>
</template>
<style lang="less" scoped>
	.form {
		width: 50%;
		margin-left: auto;
		margin-right: auto;
		.loadedPkgs {
		    max-height: 400px;
		    overflow: auto;
        }
        .list-group {
            position: absolute;
            left: 15px;
            z-index: 1;
            max-height: 247px;
            overflow: auto;
        }
		p {
			display: inline-block;
			border: 1px #b7b7b7 solid;
    		border-radius: 8px;
			padding: 5px 10px;
			margin: 5px 15px;
			font-size: 21px;
		    line-height: 1;
		}
		.close {
			padding-left: 10px;
			float: none;
		}
		.form-control {
			margin: 10px 15px;
			float: none;
		}
		.input-checkbox {
			margin: 10px 15px;
		    height: 20px;
		    width: 20px;
		}
		.btn {
            width: 72px;
			margin: 10px 15px;
		}
	}
</style>
<script>
	export default {
		name: 'config',
		data() {
			return {
				scopes: [],
				loadedPkgs: [],
                newScope: '',
				verbose: false,
                recommendPkgs: [],
                currentSelect: 0,
                state: 'reload',
			};
		},
		mounted() {
	        window.Bus.$on('change-body', (name) => {
	            if (name == 'guru-web-config') {
	            	this.$http.get('/config', {}).then((response) => {
	            		const jsonResp = response.data;
	            		this.scopes = jsonResp.inputScopes;
	            		this.loadedPkgs = jsonResp.loadedPkgs;
	            		this.verbose = jsonResp.verbose;
	                }, (response) => {
	                    window.Bus.$emit('show-alert', response.data);
	                });
	            }
	        });
	    },
	    watch: {
	    	newScope(val) {
	    		console.log('input' + val)
                if (val.length > 1) {
                    this.$http.get('/recommend-pkgs?query=' + encodeURIComponent(val), {}).then((response) => {
                        const jsonResp = response.data;
                        console.log(jsonResp)
                        this.recommendPkgs = [];
                        for (let pkg in jsonResp) {
                            this.recommendPkgs.push({item: jsonResp[pkg], select: false});
                        }
                        this.currentSelect = this.recommendPkgs.length + 1;//no select
                    }, (response) => {
                        window.Bus.$emit('show-alert', response.data);
                    });
                }
	    	}
	    },
	    methods: {
	    	removeItem(index) {
	    		this.scopes.splice(index, 1);
	    	},
            reload() {
                this.state = 'loading';
	    	    console.log('reload', this.scopes, this.verbose)
                this.$http.post('/config', {
                    scope: this.scopes.join(','),
                    verbose: this.verbose
                }, {emulateJSON: true}).then((response) => {
                    this.state = 'reload';
                    const jsonResp = response.data;
                    this.scopes = jsonResp.inputScopes;
                    this.loadedPkgs = jsonResp.loadedPkgs;
                    this.verbose = jsonResp.verbose;
                }, (response) => {
                    this.state = 'reload';
                    window.Bus.$emit('show-alert', response.data);
                });
            },
            addScope() {
	    	    if (this.currentSelect > this.recommendPkgs.length) {
                    this.scopes.push(this.newScope);
                } else {
	    	        this.scopes.push(this.recommendPkgs[this.currentSelect].item)
                }
                this.newScope = '';
                this.recommendPkgs = [];
            },
            recommendPkgsSelect(next) {
	    	    console.log(this.currentSelect)
                const preSelect = this.currentSelect;
                if (this.currentSelect > this.recommendPkgs.length) {
                    if (next > 0) {
                        this.currentSelect = 0;
                    } else {
                        this.currentSelect = this.recommendPkgs.length - 1;
                    }
                } else {
                    this.recommendPkgs[this.currentSelect].select = false;
                    console.log(this.recommendPkgs[this.currentSelect]);
                    this.currentSelect += next;
                    if (this.currentSelect >= this.recommendPkgs.length) {
                        this.currentSelect = 0;
                    } else if (this.currentSelect < 0) {
                        this.currentSelect = this.recommendPkgs.length - 1;
                    }
                }
                this.recommendPkgs[this.currentSelect].select = true;
                console.log(this.recommendPkgs[this.currentSelect]);
                console.log(this.currentSelect)
                const rp = $('#recommend-pkgs');
                const currentScrollTop = rp.scrollTop();
                let goto = 0;
                if (this.currentSelect > 5) {
                    if (preSelect < this.recommendPkgs.length
                        && preSelect != 0) {
                        goto = 41 * (this.currentSelect - 5);
                    } else {
                        goto = 41 * (this.recommendPkgs.length - 6);
                        console.log(goto)
                    }
                }
                rp.scrollTop(goto);
            }
	    }
	}
</script>