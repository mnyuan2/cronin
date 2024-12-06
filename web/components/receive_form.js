var MyReceiveForm = Vue.extend({
    template: `<div class="pipeline-form">
        <el-form :model="form" size="small">
            <el-form-item label="名称*" label-width="76px">
                <el-input v-model="form.name"></el-input>
            </el-form-item>
            <el-form-item label="别名" label-width="76px">
                <el-input v-model="form.alias" placeholder="webhook别名，为 字母、数字、-_ 的组合"></el-input>
            </el-form-item>
            <el-form-item label="接收模板*" label-width="76px">
                <div class="input-box">
                    <pre style="min-height: 200px;"><code>{{form.receive_tmpl}}</code></pre>
                    <span class="input-header">
                        <i class="el-icon-edit" @click="tmpl_box.show = true"></i>
                    </span>
                </div>
            </el-form-item>
            
            <el-form-item label="任务*" label-width="76px">
                <div><el-button type="text" @click="configSelectBox('show')">添加<i class="el-icon-plus"></i></el-button></div>
                <div id="config-selected-box" class="sort-drag-box">
                    <div class="input-box" v-for="(conf,conf_index) in form.rule_config" @mouseover="configDetailPanel(conf, true)" @mouseout="configDetailPanel(conf, false)">
                        <div class="drag">
                            <i class="el-icon-more-outline" style="transform: rotate(90deg);"></i>
                        </div>
                        <div>
                            <b class="b">{{conf.config.type_name}}</b>-<b>{{conf.config.name}}</b>-<b class="b">{{conf.config.protocol_name}}</b>
                            <el-button :type="statusTypeName(conf.config.status)" size="mini" plain round disabled>{{conf.config.status_name}}</el-button>
                            <el-divider direction="vertical"></el-divider>
                            (<el-tooltip v-for="(var_p,var_i) in conf.config.var_fields" effect="dark" :content="var_p.remark" placement="top-start">
                                <code v-if="var_p.key != ''" style="padding: 0 2px;margin: 0 4px;cursor:pointer;color: #445368;background: #f9fdff;position: relative;"><span style="position: absolute;left: -6px;bottom: -2px;">{{var_i>0 ? ',': ''}}</span>{{var_p.key}}<span class="info-2">={{var_p.value}}</span></code>
                            </el-tooltip>)
                            <span style="margin-left: 4px;" v-show="conf.view_panel">
                                <i class="el-icon-view hover" @click="configDetailBox(conf.config)" title="查看任务信息"></i>
                                <i class="el-icon-setting hover" @click="ruleConfigBox(conf_index)" title="设置关联匹配"></i>
                            </span>
                        </div>
                        <el-row style="margin:3px 0 0 8px">
                            <el-col :span="12">关联：<b class="b" v-for="re in conf.rule">{{re.key}}:{{re.value}}</b></el-col>
                            <el-col :span="12">入参：<b class="b" v-for="pm in conf.param">{{pm.key}}:{{pm.value}}</b></el-col>
                        </el-row>
                        <span class="input-header">
                            <i class="el-icon-close" @click="removeAt(conf_index)" title="移除任务"></i>
                        </span>
                    </div>
                </div>
            </el-form-item>
            
            <el-form-item label="执行间隔" label-width="76px">
                <el-input type="number" v-model="form.interval" placeholder="单个任务完成后，等待间隔时间">
                    <span slot="append">s/秒</span>
                </el-input>
            </el-form-item>
            <el-form-item label="任务停用" label-width="76px">
                <el-tooltip class="item" effect="dark" content="存在停止、错误状态任务时流水线整体停止" placement="top-start">
                    <el-radio v-model="form.config_disable_action" label="1">停止</el-radio>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="跳过停用、错误状态任务" placement="top-start">
                    <el-radio v-model="form.config_disable_action" label="2">跳过</el-radio>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="执行停用、错误状态任务" placement="top-start">
                    <el-radio v-model="form.config_disable_action" label="3">执行</el-radio>
                </el-tooltip>
            </el-form-item>
            <el-form-item label="任务错误" label-width="76px">
                <el-tooltip class="item" effect="dark" content="任务结果错误时停止流水线" placement="top-start">
                    <el-radio v-model="form.config_err_action" label="1">停止</el-radio>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="任务结果错误时跳过继续执行" placement="top-start">
                    <el-radio v-model="form.config_err_action" label="2">跳过</el-radio>
                </el-tooltip>
            </el-form-item>
            
            <el-form-item label="备注"  label-width="76px">
                <el-input v-model="form.remark"></el-input>
            </el-form-item>
            <el-form-item  label-width="76px">
                <div><el-button type="text" @click="msgBoxShow(-1)">推送<i class="el-icon-plus"></i></el-button></div>
                <div class="input-box" v-for="(msg,msg_index) in form.msg_set">
                    <el-row v-html="msg.descrition"></el-row>
                    <span class="input-header">
                        <i class="el-icon-close" @click="msgSetDel(msg_index)"></i>
                        <i class="el-icon-edit" @click="msgBoxShow(msg_index,msg)"></i>
                    </span>
                    
                </div>
            </el-form-item>
        </el-form>
        <div class="el-dialog__footer">
            <el-button size="small" @click="close()">取 消</el-button>
            <el-button type="primary" size="small" @click="submitForm()" v-if="(form.status==Enum.StatusDisable || form.status==Enum.StatusFinish || form.status==Enum.StatusError || form.status==Enum.StatusReject) && $auth_tag.pipeline_set">保存草稿</el-button>
        </div>
        
        <!-- 模板编辑弹窗 -->
        <el-dialog title="接收模板" :visible.sync="tmpl_box.show" top="10vh" class="config-select-wrap" :modal="false" :show-close="false">
            <span slot="title">接收模板
                <el-tooltip effect="dark" content="实现的参数会传入设置过参数的任务中，点击查看更多" placement="top-start">
                    <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                </el-tooltip>
            </span>
            <el-row>
                <el-input type="textarea" v-model="form.receive_tmpl" :rows="10" placeholder="接收解析模板，需返回指定json字符串结果"></el-input>
            </el-row>
            <el-row style="text-align: center;padding: 14px 0px 6px;">
                <el-button size="small" type="primary" @click="tmpl_box.show = false">确定</el-button>
            </el-row>
            <el-row style="max-height:500px;overflow-y: auto;">
                <p class="h4">模板响应样例</p>
                <div class="input-box">
                    <pre><code>{{tmpl_box.help.json}}</code></pre>
                </div>
                <el-table :data="tmpl_box.help.table" style="width: 100%">
                    <el-table-column prop="field" label="字段"></el-table-column>
                    <el-table-column prop="type" label="类型" ></el-table-column>
                    <el-table-column prop="required" label="必填"></el-table-column>
                    <el-table-column prop="desc" label="描述"></el-table-column>
                </el-table>
            </el-row>
        </el-dialog>
        
        <!-- 任务选择弹窗 -->
        <el-dialog title="任务选择" :visible.sync="config.boxShow" width="60%" top="10vh" class="config-select-wrap" :modal="false">
            <my-config-select v-if="config.boxShow" ref="selection"></my-config-select>
            <div slot="footer" class="dialog-footer">
                <a href="/index#/config" target="_blank" class="el-button el-button--text left">管理任务</a>
                <el-button size="medium" @click="configSelectBox('close')">关闭</el-button>
                <el-button size="medium" type="primary" @click="configSelectBox('confirm')" :disabled="config.running">添加</el-button>
            </div>
        </el-dialog>
        
        <!-- 推送设置弹窗 -->
        <el-dialog title="推送设置" :visible.sync="msg_set_box.show" :show-close="false" :close-on-click-modal="false" :modal="false">
            <el-form :model="msg_set_box.form" :inline="true" size="small">
                <el-form-item label="当执行">
                    <el-select v-model="msg_set_box.form.status" multiple style="width: 143px" placeholder="状态">
                        <el-option v-for="(dic_v,dic_k) in msg_set_box.statusList" :label="dic_v.name" :value="dic_v.id"></el-option>
                    </el-select>
                    时
                </el-form-item>
                <el-form-item label="发送">
                    <el-select v-model="msg_set_box.form.msg_id" style="width: 180px" placeholder="模板">
                        <el-option v-for="(dic_v,dic_k) in dic.msg" :label="dic_v.name" :value="dic_v.id"></el-option>
                    </el-select>
                    消息
                </el-form-item>
                <el-form-item label="并且@">
                    <el-select v-model="msg_set_box.form.notify_user_ids" multiple style="width: 210px" placeholder="人员">
                        <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="msg_set_box.show = false" size="small">取 消</el-button>
                <el-button type="primary" @click="msgSetConfirm()" size="small">确 定</el-button>
            </span>
        </el-dialog>
        <el-dialog title="任务详情" :visible.sync="config_detail.show" :close-on-click-modal="false" class="config-form-box" :modal="false">
            <my-config-form v-if="config_detail.show" :request="{detail:config_detail.detail,disabled:true}" @close="configDetailBox()"></my-config-form>
        </el-dialog>
        
        <el-dialog title="任务关联设置" :visible.sync="rule_config_box.show" :close-on-click-modal="false" class="config-form-box" :modal="false" width="500px">
            <el-form :model="rule_config_box.form" size="small">
                <el-form-item label="关联匹配">
                    <el-input class="input-input" v-for="(rule_v,rule_i) in rule_config_box.form.rule" v-model="rule_v.value" placeholder="匹配值">
                        <el-select v-model="rule_v.key" placeholder="匹配字段" slot="prepend" @input="e=>inputChangeArrayPush(e,rule_i,rule_config_box.form.rule)">
                            <el-option v-for="fd in dic.field" :label="fd.key" :value="fd.key">
                                <span style="float: left">{{fd.key}}</span>
                                <span style="float: right; color: #8492a6; font-size: 13px">{{fd.name}}</span>
                            </el-option>
                        </el-select>
                        <el-button slot="append" icon="el-icon-delete" @click="arrayDelete(rule_i, rule_config_box.form.rule)"></el-button>
                    </el-input>
                </el-form-item>
                <el-form-item label="入参匹配">
                    <br/>
                    <el-row v-for="(param_v,param_i) in rule_config_box.form.param">
                        <el-input v-model="param_v.key" disabled style="width: 146px;"></el-input>
                        <el-select v-model="param_v.value" placeholder="匹配字段" style="width: 280px;">
                            <el-option v-for="fd in dic.field" :label="fd.key" :value="fd.key">
                                <span style="float: left">{{fd.key}}</span>
                                <span style="float: right; color: #8492a6; font-size: 13px">{{fd.name}}</span>
                            </el-option>
                        </el-select>
                        <el-tooltip effect="dark" :content="param_v.remark" placement="top-start">
                            <i class="el-icon-info info-2"></i>
                        </el-tooltip>
                    </el-row>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button size="small" @click="ruleConfigBox(-1)">取 消</el-button>
                <el-button type="primary" size="small" @click="ruleConfigConfirm">确 定</el-button>
            </span>
        </el-dialog>
    </div>`,

    name: "MyReceiveForm",
    props: {
        request:{
            detail:Object // 详情对象
        }
    },
    data() {
        return {
            sys_info: {},
            dic:{
                user: [],
                msg: [],
                field: [],
            },
            // 表单
            form:{},
            // 接收模板弹窗
            tmpl_box:{
                show:false,
                help:{
                    json:"{\n" +
                        "    \"title\": \"xxx操作名称\",\n" +
                        "    \"user\": \"用户名\",\n" +
                        "    \"dataset\": [\n" +
                        "        {\"type\":\"jenkins\", \"repo\":\"B\", \"service\":\"s1\"},\n" +
                        "        {\"type\":\"pr\", \"repo\":\"B\", \"owner\":\"mnyuan\", \"number\":\"123\"}\n" +
                        "    ]\n" +
                        "}",
                    table:[
                        {field:'title', type:'string', required:'否',desc:'请求的标题'},
                        {field:'user', type:'string', required:'否',desc:'请求用户名称'},
                        {field:'dataset', type:'array', required:'是',desc:'匹配数据集'},
                        {field:'dataset.owner', type:'string', required:'否',desc:'空间'},
                        {field:'dataset.repo', type:'string', required:'否',desc:'仓库'},
                        {field:'dataset.number', type:'string', required:'否',desc:'序号'},
                        {field:'dataset.type', type:'string', required:'否',desc:'类型'},
                        {field:'dataset.service', type:'string', required:'否',desc:'服务'},
                        {field:'dataset.custom_1', type:'string', required:'否',desc:'自定义1'},
                        {field:'dataset.custom_2', type:'string', required:'否',desc:'自定义2'},
                    ],
                },
            },
            // 任务弹窗
            config:{
                boxShow:false,
                running: false,// 执行中
            },
            // 任务详情
            config_detail:{
                show:false,
                detail:{}
            },
            // 任务关联规则
            rule_config_box:{
                show:false,
                index: -1,
                form:{}
            },
            // 消息设置弹窗
            msg_set_box:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                form: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"结束"}, {id:0,name:"开始"}],
            },
            preference:{
                pipeline: {}
            },
            sort: null,
        }
    },
    created(){
        this.getDic()
        this.getPreference()
    },
    // 模块初始化
    mounted(){
        if (this.request.detail !== undefined && this.request.detail.id > 0){
            this.form = this.editData(this.request.detail)
        }else{
            this.form = this.addData()
        }
        this.configSort()
    },
    // 具体方法
    methods:{
        addData(){
            return  {
                id: 0,
                name: '',
                alias: '',
                type: '2',
                config_ids:[], // 任务id集合
                rule_config:[], // 任务集合
                config_disable_action: (this.preference.pipeline.config_disable_action ?? 1).toString(),
                config_err_action: '1',
                interval: this.preference.pipeline.interval ?? 0,
                msg_set: [],
                status: Enum.StatusDisable,
            }
        },
        // 编辑弹窗
        editData(row){
            row = copyJSON(row)
            for (let i in row.msg_set){
                row.msg_set[i] = this.msgSetBuildDesc(row.msg_set[i])
            }
            row.rule_config.forEach(function (item) {
                item.view_panel = false
            })
            row.config_disable_action = row.config_disable_action.toString()
            row.config_err_action = row.config_err_action.toString()
            return row
        },
        // 表单提交
        submitForm(){
            if (this.form.name == ''){
                return this.$message.error('请输入任务名称')
            }else if (this.form.spec == ''){
                return this.$message.error('请输入任务执行时间')
            }

            // 主要是强制类型
            let data = copyJSON(this.form)
            let body = {
                id: data.id,
                name: data.name,
                alias: data.alias,
                type: Number(data.type),
                spec: data.spec,
                receive_tmpl: data.receive_tmpl,
                config_ids: [],
                rule_config:data.rule_config,
                remark: data.remark,
                msg_set: data.msg_set,
                interval: Number(data.interval),
                config_disable_action: Number(data.config_disable_action),
                config_err_action: Number(data.config_err_action),
            }
            data.rule_config.forEach(function (item,index) {
                body.config_ids.push(item.config.id)
            })
            if (body.interval < 0){
                body.interval = 0
            }

            api.innerPost("/receive/set", body, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.close(true)
            })
        },


        // 推送弹窗
        msgBoxShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (oldData == undefined || index < 0){
                oldData = {
                    status: [],
                    msg_id: "",
                    notify_user_ids: [],
                }
            }else if (typeof oldData != 'object'){
                console.log('推送信息异常', oldData)
                return this.$message.error("推送信息异常");
            }
            this.msg_set_box.show = true
            this.msg_set_box.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.msg_set_box.title = this.msg_set_box.index < 0? '添加' : '编辑';
            this.msg_set_box.form = copyJSON(oldData)
        },
        msgSetDel(index){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('确认删除推送配置','提示',{type:'warning'}).then(()=>{
                this.form.msg_set.splice(index, 1);
            })
        },
        // 推送确认
        msgSetConfirm(){
            if (this.msg_set_box.form.msg_id <= 0){
                return this.$message.warning("请选择消息模板");
            }
            let data = this.msgSetBuildDesc(this.msg_set_box.form)

            if (this.msg_set_box.index < 0){
                this.form.msg_set.push(data)
            }else{
                this.form.msg_set[this.msg_set_box.index] = data
            }
            this.msg_set_box.show = false
            this.msg_set_box.index = -1
            this.msg_set_box.form = {}
        },
        // 构建消息设置描述
        msgSetBuildDesc(data){
            let item1 = this.msg_set_box.statusList.filter((option) => {
                return data.status.includes(option.id)
            }).map((item)=>{return item.name});
            data.status_name = item1.join(',')

            let descrition = '<i class="el-icon-bell"></i>当任务<b class="b">'+data.status_name+'</b>时'

            let item2 = this.dic.msg.find(option => option.id === data.msg_id)
            if (item2){
                data.msg_name = item2.name
                descrition += '，发送<b class="b">'+item2.name+'</b>消息'
            }
            let item3 = this.dic.user.filter((option) => {
                return data.notify_user_ids.includes(option.id);
            }).map((item)=>{return item.name})
            if (item3.length > 0){
                data.notify_users_name = item3.join(',')
                descrition += '，并且@人员<b class="b">'+data.notify_users_name+'</b>'
            }
            data.descrition = descrition
            return data
        },
        // 枚举
        getDic(){
            api.dicList([Enum.dicUser, Enum.dicMsg, Enum.dicReceiveDataField],(res) =>{
                this.dic.user = res[Enum.dicUser]
                this.dic.msg = res[Enum.dicMsg]
                this.dic.field = res[Enum.dicReceiveDataField]
            })
        },

        // 任务盒子弹窗
        configSelectBox(e='show'){
            if (e == 'show'){ // 显示
                this.config.boxShow = true
            }else if (e == 'close'){ // 关闭
                this.config.boxShow = false
            }else if (e == 'confirm'){ // 提交
                console.log("选中元素",)
                this.config.running = true
                this.$refs.selection.selected.forEach((item)=>{
                    this.form.rule_config.push({rule:[],param:[],config:copyJSON(item), view_panel: false})
                })
                this.config.boxShow = false
                this.config.running = false
            }
            console.log("任务盒子",this.config)
        },
        configSort(){
            const that = this
            this.$nextTick(() => {
                that.sort = MySortable(document.getElementById("config-selected-box"), (oldIndex, newIndex)=>{
                    const oldlist = copyJSON(that.form.rule_config);
                    const oldItem = oldlist.splice(oldIndex, 1)[0];
                    oldlist.splice(newIndex, 0, oldItem);

                    that.form.rule_config = []
                    that.$nextTick((t)=>{
                        this.form.rule_config= oldlist
                        console.log("拖拽后",that.form.rule_config, this, t)
                    })
                })
            })
        },
        configDetailPanel(detail,show=false){
            detail.view_panel = show === true
        },
        configDetailBox(detail=null){
            if (!detail){
                this.config_detail.show =false
                this.config_detail.detail = {}
                return
            }
            api.innerGet('/config/detail', {id: detail.id}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.config_detail.show = true
                this.config_detail.detail = res.data
            },{async:false})

        },
        removeAt(idx) {
            this.$confirm('确定移除任务', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.form.rule_config.splice(idx, 1);
            }).catch(() => {/*取消*/});
        },
        // 任务关联盒子
        ruleConfigBox(index){
            let data = copyJSON(this.form.rule_config[index]??{})
            if (!data.config){
                this.rule_config_box.show = false
                this.rule_config_box.index = -1
                this.rule_config_box.form = {}
                return
            }

            if (data.rule == null || data.rule.length == 0){
                data.rule = [{key:'',value:''}]
            }else{
                data.rule.push({key:'',value:''})
            }

            // 解析配置里面的参数
            let old_param_map = {}
            data.param.forEach(function (item) {
                old_param_map[item.key] = item.value
            })
            let param = []
            data.config.var_fields.forEach(function (item) {
                // 入参key / 匹配字段 / 入参描述
                param.push({key:item.key,value:old_param_map[item.key]??'',remark:item.remark})
            })
            data.param = param

            this.rule_config_box.show = true
            this.rule_config_box.index = Number(index)
            this.rule_config_box.form = data
            console.log('rule_config_box.form', data)
        },
        // 关联设置确认
        ruleConfigConfirm(){
            let data = copyJSON(this.rule_config_box.form)
            data.rule = data.rule.filter(function (rule) {
                return rule.key !== undefined && rule.key != '' && rule.value != ''
            })
            data.param = data.param.filter(function (param) {
                return param.key !== undefined && param.key != '' && param.value != ''
            })
            this.form.rule_config[this.rule_config_box.index] = data
            this.rule_config_box.show = false
        },
        // 获取偏好
        getPreference(){
            api.innerGet("/setting/preference_get", null,res=>{
                if (!res.status){
                    return this.$message.error('偏好错误，'+res.message)
                }
                this.preference = res.data
            },{async: false})
        },
        close(is_change=false){
            this.$emit('close', {is_change:is_change})
        },
    }
})

Vue.component("MyReceiveForm", MyReceiveForm);