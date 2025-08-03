var MyHeader = Vue.extend({
    template: `<el-header height="46px">
                        <el-menu router
                          :default-active="$route.path"
                          class="menu-left"
                          mode="horizontal"
                          background-color="#151515" 
                          text-color="#fff"
                          active-text-color="#66b1ff">
                          <el-menu-item index="/" style="font-weight: 500;font-size: 110%;">cronin</el-menu-item>
                          <el-menu-item index="/config" v-if="$auth_tag.config_list">任务</el-menu-item>
                          <el-menu-item index="/pipeline" v-if="$auth_tag.pipeline_list">流水线</el-menu-item>
                          <el-menu-item index="/receive" v-if="$auth_tag.receive_list">接收</el-menu-item>
                          <el-menu-item index="/logs" >日志</el-menu-item> <!-- v-if="$auth_tag.logs_list"   -->
                          <!--右导航-->
                          <el-menu-item-group class="group-right">
                              <el-menu-item>
                                <el-dropdown @command="envClick">
                                    <span class="el-dropdown-link">
                                        {{sys_info.env_name}} <i class="el-icon-arrow-down el-icon--right"></i>
                                    </span>
                                    <el-dropdown-menu slot="dropdown">
                                        <el-dropdown-item v-for="(dic_v,dic_k) in dic_envs" :command="dic_v.key" :disabled="sys_info.env==dic_v.key">{{dic_v.name}}</el-dropdown-item>
                                        <el-dropdown-item command="envBoxDisplay" divided v-if="$auth_tag.env_list">管理环境</el-dropdown-item>
                                    </el-dropdown-menu>
                                </el-dropdown>
                              </el-menu-item>
                              <el-menu-item index="/my/work" class="item" title="我的待办工作">工作</el-menu-item>
                              <el-submenu popper-class="submenu" index="user" class="item">
                                <template slot="title">{{user.username}}<i class="el-submenu__icon-arrow el-icon-arrow-down"></template>
                                <el-menu-item :index="'/user?data_id='+user.id">账号信息</el-menu-item>
                                <el-divider style="margin: 0"></el-divider>
                                <el-menu-item @click="logout">退出登录</el-menu-item>
                              </el-submenu>
                              <el-submenu popper-class="submenu" index="setting" class="icon">
                                <template slot="title"><i class="el-icon-more" title="设置及其他"></template>
                                <el-menu-item index="/message_template" v-if="$auth_tag.message_list">通知模板</el-menu-item>
                                <el-menu-item index="/tag" v-if="$auth_tag.tag_list">标签</el-menu-item>
                                <el-menu-item index="/users" v-if="$auth_tag.user_list">人员</el-menu-item>
                                <el-menu-item index="/role" v-if="$auth_tag.role_list">权限</el-menu-item>
                                <el-menu-item index="/source" v-if="$auth_tag.source_list">连接</el-menu-item>
                                <el-menu-item index="/setting">设置</el-menu-item>
                                <el-divider></el-divider>
                                <el-menu-item><a href="https://cron.qqe2.com/" target="_blank">时间格式生成器</a></el-menu-item>
                                <el-menu-item><a href="https://gitee.com/mnyuan/cronin/" target="_blank">Gitee <i class="el-icon-star-off" style="vertical-align: initial;font-size: 15px;"></i></a></el-menu-item>
                                <el-menu-item><a href="https://gitee.com/mnyuan/cronin/issues/new" target="_blank">反馈</a></el-menu-item>
                                <el-menu-item disabled>cronin {{sys_info.version}}</el-menu-item>
                              </el-submenu>
                          </el-menu-item-group>
                        </el-menu>
                    
                    
                    <!--
                        <el-menu router 
                            :default-active="$route.path"
                            class="menu-right" 
                            mode="horizontal" 
                            background-color="#151515"  
                            text-color="#fff" 
                            active-text-color="#409effa8">
                          <el-menu-item>
                            <el-dropdown @command="envClick">
                                <span class="el-dropdown-link">
                                    {{sys_info.env_name}} <i class="el-icon-arrow-down el-icon--right"></i>
                                </span>
                                <el-dropdown-menu slot="dropdown">
                                    <el-dropdown-item v-for="(dic_v,dic_k) in dic_envs" :command="dic_v.key" :disabled="sys_info.env==dic_v.key">{{dic_v.name}}</el-dropdown-item>
                                    <el-dropdown-item command="envBoxDisplay" divided>管理环境</el-dropdown-item>
                                </el-dropdown-menu>
                            </el-dropdown>
                          </el-menu-item>
                          <el-submenu popper-class="submenu">
                            <template slot="title">设置</template>
                            <el-menu-item index="/message_template">通知</el-menu-item>
                            <el-menu-item index="/users">人员</el-menu-item>
                          </el-submenu>
                          <el-submenu popper-class="submenu">
                            <template slot="title">关于</template>
                            <el-menu-item><a href="https://cron.qqe2.com/" target="_blank">时间格式生成器</a></el-menu-item>
                            <el-menu-item><a href="https://gitee.com/mnyuan/cronin/" target="_blank">Gitee <i class="el-icon-star-off" style="vertical-align: initial;font-size: 15px;"></i></a></el-menu-item>
                            <el-menu-item disabled>cronin {{sys_info.version}}</el-menu-item>
                          </el-submenu>  
                        </el-menu>
                        -->
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
            activeIndex: 1,
            user:{},
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        this.user = cache.getUser()
        if (!this.user.id){ // 未登录
            backLoginPage()
        }
    },

    // 具体方法
    methods:{
        handleSelect(key, keyPath) {
            console.log(key, keyPath);
        },
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
                api.dicList([Enum.dicEnv], (list) => {
                    this.dic_envs = list[Enum.dicEnv]
                }, true)
            }
        },
        logout(){
            cache.empty()
            backLoginPage()
        }
    }
})

Vue.component("MyHeader", MyHeader);