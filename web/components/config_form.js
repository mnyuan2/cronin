var MyConfigForm = Vue.extend({
    template: `<div class="config-form">
    <el-form :model="form">
        <el-form-item label="活动名称*" label-width="76px">
            <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="类型*" label-width="76px" v-if="mode=='config'">
            <el-radio v-model="form.type" label="1">周期</el-radio>
            <el-radio v-model="form.type" label="2">单次</el-radio>
        </el-form-item>

        <el-form-item label="时间*" label-width="76px" v-if="mode=='config'">
            <el-input v-show="form.type==1" v-model="form.spec" :placeholder="hintSpec"></el-input>
            <el-date-picker 
                style="width: 100%"
                v-show="form.type==2" 
                v-model="form.spec" 
                value-format="yyyy-MM-dd HH:mm:ss"
                type="datetime" 
                placeholder="选择运行时间" 
                :picker-options="pickerOptions">
            </el-date-picker>
        </el-form-item>
        
        <el-form-item>
            <el-tabs type="border-card" v-model="form.protocol">
                <el-tab-pane label="http" name="1">
                    <el-form-item label="请求地址" class="http_url_box">
                        <el-input v-model="form.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                            <el-select v-model="form.command.http.method" placeholder="请选请求方式" slot="prepend">
                                <el-option label="GET" value="GET"></el-option>
                                <el-option label="POST" value="POST"></el-option>
                            </el-select>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="请求Header" class="http_header_box">
                        <el-input v-for="(header_v,header_i) in form.command.http.header" v-model="header_v.value" placeholder="参数值">
                            <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                            <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="请求Body参数">
                        <el-input type="textarea" v-model="form.command.http.body" rows="5" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                    </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="rpc" name="2">
                    <el-form-item label="请求模式">
                        <el-radio v-model="form.command.rpc.method" label="GRPC">GRPC</el-radio>
                    </el-form-item>
                    <el-form-item label="proto">
                        <el-input type="textarea" v-model="form.command.rpc.proto" rows="5" placeholder="请输入*.proto 文件内容"></el-input>
                    </el-form-item>
                    <el-form-item label="地址" label-width="42px" style="margin: 22px 0;">
                        <el-input v-model="form.command.rpc.addr" placeholder="请输入服务地址，含端口; 示例：localhost:21014"></el-input>
                    </el-form-item>
                    <el-form-item label="方法">
                        <el-select v-model="form.command.rpc.action" filterable clearable placeholder="填写proto信息后点击解析，可获得其方法进行选择" style="min-width: 400px">
                            <el-option v-for="item in form.command.rpc.actions" :key="item" :label="item" :value="item"></el-option>
                        </el-select>
                        <el-button @click="parseProto()">解析proto</el-button>
                    </el-form-item>
                    <el-form-item label="请求参数">
                        <el-input type="textarea" v-model="form.command.rpc.body" rows="5" placeholder="请输入请求参数"></el-input>
                    </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="cmd" name="3">
                    <el-input type="textarea" v-model="form.command.cmd" rows="5" :placeholder="sys_info.cmd_name + ': 请输入命令行执行内容'"></el-input>
                </el-tab-pane>
                <el-tab-pane label="sql" name="4" label-position="left">
                    <el-form-item label="驱动">
                        <el-select v-model="form.command.sql.driver" placeholder="驱动">
                            <el-option label="mysql" value="mysql"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="链接">
                        <el-select v-model="form.command.sql.source.id" placement="请选择sql链接">
                            <el-option v-for="(dic_v,dic_k) in dic_list.sql_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                        <el-button type="text" style="margin-left: 20px" @click="sourceBox(true)">设置链接</el-button>
                    </el-form-item>
                    <el-form-item label="执行语句">
                        <div><el-button type="text" @click="sqlSetShow(-1,'')">添加<i class="el-icon-plus"></i></el-button></div>
                        <div style="overflow-y: auto;max-height: 420px;">
                            <div v-for="(statement,sql_index) in form.command.sql.statement" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                <pre style="margin: 0;overflow-y: auto;max-height: 180px;min-height: 56px;"><code class="language-sql hljs">{{statement}}</code></pre>
                                <i class="el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="sqlSetDel(sql_index)"></i>
                                <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="sqlSetShow(sql_index,statement)"></i>
                                <i style="position: absolute;right: 1px;top: 49px;font-size: 16px;">#{{sql_index}}</i>
                            </div>
                            <el-alert v-show="form.command.sql.statement.length==0" title="未添加执行sql，请添加。" type="info"></el-alert>
                        </div>
                    </el-form-item>
                    <el-form-item label="错误行为">
                        <el-tooltip class="item" effect="dark" content="执行到错误停止，之前的成功语句保留" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="1">终止任务</el-radio>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="执行到错误跳过，继续执行后面语句" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="2">跳过继续</el-radio>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="执行到错误，之前的成功语句也整体回滚" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="3">事务回滚</el-radio>
                        </el-tooltip>
                    </el-form-item>
                    <el-form-item label="执行间隔">
                        <el-input v-model="form.command.sql.interval" oninput ="value=value.replace(/[^\d]/g,'')">
                            <template slot="append">秒</template>
                        </el-input>
                    </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="jenkins" name="5">
                    <el-form-item label="链接">
                        <el-select v-model="form.command.jenkins.source.id" placement="请选择链接">
                            <el-option v-for="(dic_v,dic_k) in dic_list.jenkins_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                        <el-button type="text" style="margin-left: 20px" @click="sourceBox(true)">设置链接</el-button>
                    </el-form-item>
                    <el-form-item label="项目">
                        <el-input v-model="form.name"></el-input>
                    </el-form-item>
                    <el-form-item label="请求Header" class="http_header_box">
                        <el-input v-for="(header_v,header_i) in form.command.http.header" v-model="header_v.value" placeholder="参数值">
                            <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                            <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                        </el-input>
                    </el-form-item>
                </el-tab-pane>
                
            </el-tabs>
        </el-form-item>
        <el-form-item label="备注" label-width="43px">
            <el-input v-model="form.remark"></el-input>
        </el-form-item>
        <el-form-item label-width="2px">
            <div><el-button type="text" @click="msgBoxShow(-1)">推送<i class="el-icon-plus"></i></el-button></div>
            <div v-for="(msg,msg_index) in form.msg_set" style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                <el-row v-html="msg.descrition"></el-row>
                <i class="el-tag__close el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="msgSetDel(msg_index)"></i>
                <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="msgBoxShow(msg_index,msg)"></i>
            </div>
        </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
        <el-button @click="configRun()" class="left" v-show="form.type==1" v-if="mode=='config'">执行一下</el-button>
        <el-button @click="setConfigShow = false">取 消</el-button>
        <el-button type="primary" @click="setCron()">确 定</el-button>
    </div>
    
    <!-- sql链接源管理弹窗 -->
    <el-drawer title="链接管理" :visible.sync="source.boxShow" size="40%" wrapperClosable="false" :before-close="sourceBox(-1)">
        <my-sql-source></my-sql-source>
    </el-drawer>
</div>`,
    name: "MyConfigForm",
    props: {
        config_id:Number, // 任务id，0.表示新增
        mode:String, // 模式：config.常规任务(默认)、pipeline.流水线
    },
    data(){
        return {
            config_id:0,
            mode: 'config',

            sys_info:{},
            dic_list:{
                sql_source:[],
                jenkins_source:[],
                user:[],
                msg:[],
            },

            form:{},
            hintSpec: "* * * * * *",
            source: {
                boxShow: false,
                dic_type: 0,
            },
            sqlSet: {
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: "", // 实际内容
            }, // sql设置弹窗
            msgSet:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"成功"}, {id:0,name:"完成"}],
            },
            // 日期选择器设置
            pickerOptions: {
                disabledDate(time){
                    return time.getTime() < Date.now() - 8.64e7
                },
                selectableRange: "00:00:00 - 23:01:59",
            },
        }
    },
    // 模块初始化
    created(){
        this.form = this.initFormData()
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
        this.getDicSource()
    },
    // 模块初始化
    mounted(){},
    watch:{
        config_id:{
            immediate: true, // 解决首次负值不触发的情况
            handler: function (newVal,oldVal){
                console.log("config_log config_id",newVal, oldVal)
                if (newVal != 0){
                    this.logByConfig(newVal)
                }
            },
        }
    },

    // 具体方法
    methods:{
        initFormData(){
            return  {
                type: '1',
                protocol: '3',
                command:{
                    http:{
                        method: 'GET',
                        header: [{"key":"","value":""}],
                        url:'',
                        body:'',
                    },
                    rpc:{
                        proto: '',
                        method: 'GRPC',
                        addr: '',
                        action: '',
                        actions: [],
                        header: [],
                        body: ''
                    },
                    cmd:'',
                    sql:{
                        driver: "mysql",
                        source:{
                            id: "",
                            // title: "",
                            // hostname:"",
                            // port:"",
                            // username:"",
                            // password:""
                        },
                        statement:[],
                        err_action: "1",
                        interval: 0, // 与事务回滚互斥
                    },
                    jenkins:{
                        source:{
                            id: "",
                        },
                        name: "",
                        params: [{"key":"","value":""}]
                    }
                },
                msg_set: []
            }
        },
        // http header 输入值变化时
        httpHeaderInput(val){
            if (val == ""){
                return
            }
            let item = this.form.command.http.header.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.command.http.header.push({"key":"","value":""})
            }
        },
        // http header 输入值删除
        httpHeaderDel(index){
            if ((index+1) >= this.form.command.http.header.length){
                return
            }
            this.form.command.http.header.splice(index,1)
        },
        // 枚举
        getDicSource(){
            api.dicList([Enum.dicSqlSource, Enum.dicJenkinsSource, Enum.dicUser, Enum.dicMsg],(res) =>{
                this.dic_list.sql_source = res[Enum.dicSqlSource]
                this.dic_list.jenkins_source = res[Enum.dicJenkinsSource]
                this.dic_list.user = res[Enum.dicUser]
                this.dic_list.msg = res[Enum.dicMsg]
            })
        },
        // 解析proto内容
        parseProto(){
            let data = this.form.command.rpc.proto
            if (data == "" || data == undefined){
                return this.$message.warning("请先输入proto内容");
            }
            api.innerPost("/foundation/parse_proto", {proto:data}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.form.command.rpc.actions = res.data.actions
            })
        },
        // source box
        sourceBox(e){
            if (e == -1){
                this.source.boxShow = false
                // this.getDicSource() // 关闭弹窗要重载枚举 此处会导致死循环
                return
            }else if (e > 1){
                this.source.dic_type = e
                this.source.boxShow = true
            }else{
                return this.$message.warning('参数不规范')
            }
        },
        // 推送弹窗
        msgBoxShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (oldData == undefined || index < 0){
                oldData = {
                    status: 1,
                    msg_id: "",
                    notify_user_ids: [],
                }
            }else if (typeof oldData != 'object'){
                console.log('推送信息异常', oldData)
                return this.$message.error("推送信息异常");
            }
            this.msgSet.show = true
            this.msgSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.msgSet.title = this.msgSet.index < 0? '添加' : '编辑';
            this.msgSet.data = oldData
        },
    }
})

Vue.component("MyConfigForm", MyConfigForm);