var MyHeader = Vue.extend({
    template: `<el-header>
                    <el-dropdown @command="envClick">
                        <span class="el-dropdown-link">
                            {{sys_info.env_name}} <i class="el-icon-arrow-down el-icon--right"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item v-for="(dic_v,dic_k) in dic_envs" :command="dic_v.key" :disabled="sys_info.env==dic_v.key">{{dic_v.name}}</el-dropdown-item>
                            <el-dropdown-item command="envBoxDisplay" divided>管理环境</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                    <div style="float: right">
                        <a href="https://cron.qqe2.com/" target="_blank" >6位 时间格式生成器</a>
                        <span class="version">{{sys_info.version}}</span>
                    </div>
                    
                    <!-- 环境 管理弹窗 -->
                    <el-drawer title="环境管理" :visible.sync="envBoxShow" size="50%" wrapperClosable="false" :before-close="envBox">
                        <my-env :reload_list="envBoxShow"></my-env>
                    </el-drawer>
                </el-header>`,
    name: "MyHeader",
    props: {
        dic_envs:[],
        sys_info:{},
        reload_list:false, // 重新加载列表
    },
    data(){
        return {
            sys_info: {},
            envBoxShow: false, // 环境管理弹窗
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){},

    // 具体方法
    methods:{
        envClick(cmd){
            if (cmd == ""){
                return
            }
            // 环境管理
            if (cmd == "envBoxDisplay"){
                return this.envBox(true)
            }
            // 环境切换
            if (this.sys_info.env == cmd){
                return
            }
            let env_name = ""
            for (let i in this.dic_envs){
                if (this.dic_envs[i].key == cmd){
                    env_name = this.dic_envs[i].name
                    break
                }
            }
            if (env_name == ""){
                return this.$message.error("选择环境信息异常")
            }
            // this.sys_info.env = cmd
            // this.sys_info.env_name = env_name
            // api.setEnv(cmd, env_name)
            window.location.assign(window.location.protocol+"//" + window.location.host + window.location.pathname + "?env="+cmd)
            // console.log("url",window.location.protocol, window.location.host, window.location.pathname)
        },
        envBox(show){
            this.envBoxShow = show == true;
            if (!this.envBoxShow){
                this.getDic() // 关闭弹窗要重载枚举?
            }
        },
    }
})

Vue.component("MyHeader", MyHeader);