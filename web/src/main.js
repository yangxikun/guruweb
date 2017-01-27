import Vue from 'vue'
import VueResource from 'vue-resource'
import GuruWebHead from './GuruWebHead.vue'
import GuruWebSearch from './GuruWebSearch.vue'
import GuruWebFile from './GuruWebFile.vue'
import GuruWebMenu from './GuruWebMenu.vue'
import GuruWebOut from './GuruWebOut.vue'
import GuruWebAlert from './GuruWebAlert.vue'
import GuruWebConfig from './GuruWebConfig.vue'
import GuruWebHistoryNav  from './GuruWebHistoryNav.vue'

Vue.use(VueResource);

window.Bus = new Vue;

new Vue({
    el: '#app',
    data() {
        return {
            guruWebBody: ''
        }
    },
    mounted() {
        window.Bus.$on('change-body', (name) => {
            this.guruWebBody = name
        });
    },
    components: {
        GuruWebHead,
        GuruWebSearch,
        GuruWebFile,
        GuruWebOut,
        GuruWebMenu,
        GuruWebAlert,
        GuruWebConfig,
        GuruWebHistoryNav
    }
});
