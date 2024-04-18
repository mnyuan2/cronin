var MyPipeline = Vue.extend({
    template: `<el-main>
                    <el-menu :default-active="labelType" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
<!--                        <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>-->
<!--                        <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>-->
                        <div style="float: right">
                            <el-button type="text" @click="formBox(0)">添加流水线</el-button>
<!--                            <el-button type="text" @click="getRegisterList()">已注册任务</el-button>-->
                        </div>
                    </el-menu>
                    <el-table :data="list.items">
                        <el-table-column prop="id" label="编号"></el-table-column>
                        <el-table-column prop="" label="执行成功率">
                            <template slot-scope="scope">
                                <el-tooltip placement="top-start">
                                    <div slot="content">{{scope.row.topRatio}}%<br/>近期{{scope.row.top_number}}次执行，{{scope.row.top_error_number}}次失败。</div>
                                    <i :class="getTopIcon(scope.row.top_number, scope.row.topRatio)"></i>
                                </el-tooltip>
                            </template>
                        </el-table-column>
                        <el-table-column prop="spec" label="执行时间"></el-table-column>
                        <el-table-column prop="name" label="任务名称"></el-table-column>
                        <el-table-column prop="remark" label="备注"></el-table-column>
                        <el-table-column prop="" label="状态">
                            <template slot-scope="scope">
                                <el-tooltip placement="top-start">
                                    <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                                    <span :class="statusClass(scope.row.status)">{{scope.row.status_name}}</span>
                                </el-tooltip>
                            </template>
                        </el-table-column>
                        <el-table-column label="操作">
                            <template slot-scope="{row}">
                                <el-button type="text" @click="formBox(row)">编辑</el-button>
                                <el-button type="text" @click="changeStatus(row, 2)" v-if="row.status!=2">激活</el-button>
                                <el-button type="text" @click="changeStatus(row, 1)" v-if="row.status==2">停用</el-button>
                                <el-button type="text" @click="configLogBox(row.id, row.name)">日志</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                    <el-pagination
                            @size-change="handleSizeChange"
                            @current-change="handleCurrentChange"
                            :current-page.sync="list.page.page"
                            :page-size="list.page.size"
                            layout="total, prev, pager, next"
                            :total="list.page.total">
                    </el-pagination>
                    
                    
                    <!-- 流水线设置表单 :before-close="configLogBox(0)" -->
                    <el-drawer :title="form.boxTitle" :visible.sync="form.boxShow" size="60%" wrapperClosable="false">
                        <el-form :model="form">
                            <el-form-item label="名称*" label-width="76px">
                                <el-input v-model="form.data.name"></el-input>
                            </el-form-item>
<!--                            <el-form-item label="类型*" label-width="76px">-->
<!--                                <el-radio v-model="form.type" label="1">周期</el-radio>-->
<!--                                <el-radio v-model="form.type" label="2">单次</el-radio>-->
<!--                            </el-form-item>-->
            
                            <el-form-item label="时间*" label-width="76px">
                                <el-input v-show="form.data.type==1" v-model="form.data.spec" :placeholder="form.hintSpec"></el-input>
                                <el-date-picker 
                                    style="width: 100%"
                                    v-show="form.data.type==2" 
                                    v-model="form.data.spec" 
                                    value-format="yyyy-MM-dd HH:mm:ss"
                                    type="datetime" 
                                    placeholder="选择运行时间" 
                                    :picker-options="form.pickerOptions">
                                </el-date-picker>
                            </el-form-item>
                            <el-form-item prop="plate" label-width="76px">
                                <span slot="label" style="white-space: nowrap;">
                                    参数实现
                                    <el-tooltip effect="dark" content="点击查看变量参数使用详解" placement="top-start">
                                        <i class="el-icon-info" style="cursor: pointer"></i>
                                    </el-tooltip>
                                </span>
                                <el-input type="textarea" v-model="form.var_params" placeholder="变量参数实现 json 格式"></el-input>
                            </el-form-item>
                            
                            <el-form-item label="任务" label-width="76px">
                                <div><el-button type="text" @click="configSelectBox('show')">添加<i class="el-icon-plus"></i></el-button></div>
                                <div id="config-selected-box" class="sort-drag-box">
                                    <div class="item-drag" v-for="(conf,conf_index) in form.data.configs" style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                        <i class="el-icon-more-outline drag"></i>
                                        <el-tag type="info">{{conf.type_name}}</el-tag>-
                                        <el-tag>{{conf.name}}</el-tag>-
                                        <el-tag type="info">{{conf.protocol_name}}</el-tag>
                                        <el-tag>{{conf.status_name}}</el-tag>
                                        <el-divider direction="vertical"></el-divider>
                                            <el-tooltip v-for="(var_p,var_i) in conf.var_params" effect="dark" :content="var_p.value" placement="top-start">
                                                <code v-if="var_p.key != ''" style="padding: 0 4px;margin-right: 4px;cursor:pointer;color: #445368;background: #f9fdff;"><{{var_p.key}}></code>
                                            </el-tooltip>
                                        <i class="el-icon-close item-close" @click="removeAt(conf_index)"></i>
                                    </div>
                                </div>
                            </el-form-item>
                            
                            <el-form-item label="任务停用" label-width="76px">
                                <el-tooltip class="item" effect="dark" content="存在停止、错误状态任务时流水线整体停止" placement="top-start">
                                    <el-radio v-model="form.data.config_disable_action" label="1">停止</el-radio>
                                </el-tooltip>
                                <el-tooltip class="item" effect="dark" content="跳过停用、错误状态任务" placement="top-start">
                                    <el-radio v-model="form.data.config_disable_action" label="2">跳过</el-radio>
                                </el-tooltip>
                                <el-tooltip class="item" effect="dark" content="执行停用、错误状态任务" placement="top-start">
                                    <el-radio v-model="form.data.config_disable_action" label="3">执行</el-radio>
                                </el-tooltip>
                            </el-form-item>
                            <el-form-item label="任务错误" label-width="76px">
                                <el-tooltip class="item" effect="dark" content="任务结果错误时停止流水线" placement="top-start">
                                    <el-radio v-model="form.data.config_err_action" label="1">停止</el-radio>
                                </el-tooltip>
                                <el-tooltip class="item" effect="dark" content="任务结果错误时跳过继续执行" placement="top-start">
                                    <el-radio v-model="form.data.config_err_action" label="2">跳过</el-radio>
                                </el-tooltip>
                            </el-form-item>
                            
                            <el-form-item label="备注"  label-width="76px">
                                <el-input v-model="form.data.remark"></el-input>
                            </el-form-item>
                            <el-form-item  label-width="76px">
                                <div><el-button type="text" @click="msgBoxShow(-1)">推送<i class="el-icon-plus"></i></el-button></div>
                                <div v-for="(msg,msg_index) in form.data.msg_set" style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                    <el-row v-html="msg.descrition"></el-row>
                                    <i class="el-tag__close el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="msgSetDel(msg_index)"></i>
                                    <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="msgBoxShow(msg_index,msg)"></i>
                                </div>
                            </el-form-item>
                        </el-form>
                        <div>
<!--                            <el-button @click="configRun()" class="left" v-show="form.type==1">执行一下</el-button>-->
                            <el-button size="small" @click="formBox(-1)">取 消</el-button>
                            <el-button type="primary" size="small" @click="submitForm()">确 定</el-button>
                        </div>
                    </el-drawer>
                    <!-- 任务日志弹窗 -->
                    <el-drawer :title="log.title" :visible.sync="log.show" direction="rtl" size="40%" wrapperClosable="false" > <!-- :before-close="configLogBox(0)"-->
                        <my-config-log :config_id="log.id"></my-config-log>
                    </el-drawer>
                    <!-- 注册任务列表弹窗 -->
                    <el-drawer title="已注册任务" :visible.sync="register.boxShow" direction="rtl" size="40%" wrapperClosable="false">
                        <el-table :data="register.items">
                            <el-table-column property="id" label="编号"></el-table-column>
                            <el-table-column property="spec" label="执行时间"></el-table-column>
                            <el-table-column property="update_dt" label="下一次执行"></el-table-column>
                            <el-table-column property="name" label="任务名称"></el-table-column>
                            <el-table-column property="" label="">
                                <template slot-scope="scope">
                                    <el-button type="text" @click="configLogBox(scope.row.id, scope.row.name)">日志</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </el-drawer>
            
                    <!-- 任务选择弹窗 -->
                    <el-dialog title="任务选择" :visible.sync="config.boxShow" width="60%" top="10vh" class="config-select-wrap">
                        <my-config-select ref="selection"></my-config-select>
                        <div slot="footer" class="dialog-footer">
                            <a href="/index#/config" target="_blank" class="el-button el-button--text left">管理任务</a>
                            <el-button size="medium" @click="configSelectBox('close')">关闭</el-button>
                            <el-button size="medium" type="primary" @click="configSelectBox('confirm')" :disabled="config.running">添加</el-button>
                        </div>
                    </el-dialog>
                    <!-- 推送设置弹窗 -->
                    <el-dialog title="推送设置" :visible.sync="msgSet.show" :show-close="false" :close-on-click-modal="false">
                        <el-form :model="msgSet" :inline="true" size="mini">
                            <el-form-item label="当">
                                <el-select v-model="msgSet.data.status" style="width: 90px">
                                    <el-option v-for="(dic_v,dic_k) in msgSet.statusList" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                                时
                            </el-form-item>
                            <el-form-item label="发送">
                                <el-select v-model="msgSet.data.msg_id">
                                    <el-option v-for="(dic_v,dic_k) in dic.msg" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                                消息
                            </el-form-item>
                            <el-form-item label="并且@用户">
                                <el-select v-model="msgSet.data.notify_user_ids" multiple="true">
                                    <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                            </el-form-item>
                        </el-form>
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="msgSet.show = false">取 消</el-button>
                            <el-button type="primary" @click="msgSetConfirm()">确 定</el-button>
                        </span>
                    </el-dialog>
                </el-main>`,
    name: "MyPipeline",
    data(){
        return {
            dic:{
                user: [],
                msg: [],
            },
            labelType: '2',
            // 列表数据
            list: {
                items: [],
                page: {
                    total:0,
                    page: 1,
                    size: 20,
                },
                param:{
                    type: 2,
                    page: 1,
                    size: 20,
                },
                request: false, // 请求中标志
            },
            // 已注册任务
            register:{
                items:[],
                boxShow: false,
            },
            // 表单
            form:{
                data:[],
                boxShow: false, // 新增弹窗
                boxTitle: '',
                hintSpec: "* * * * * *",
                // 日期选择器设置
                pickerOptions: {
                    disabledDate(time){
                        return time.getTime() < Date.now() - 8.64e7
                    },
                    selectableRange: "00:00:00 - 23:01:59",
                },
            },
            // 任务弹窗
            config:{
                boxShow:false,
                running: false,// 执行中
            },
            // 消息设置弹窗
            msgSet:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"结束"}, {id:0,name:"开始"}],
            },
            // 日志弹窗
            log:{
                show: false,
                id: 0,
                title:""
            },

        }
    },
    // 模块初始化
    created(){
        this.form.data = this.initFormData()
        this.getDic()
    },
    // 模块初始化
    mounted(){
        this.getList()
    },
    watch:{
        "form.spec":{
            handler(v){
                if (this.form.data.type == 2){ // 年月日的变化控制
                    let cur = moment()
                    if (moment(v).format('YYYY-MM-DD') == cur.format('YYYY-MM-DD')){
                        this.form.pickerOptions.selectableRange = `${cur.format('HH:mm:ss')} - 23:59:59`
                    }else{
                        this.form.pickerOptions.selectableRange = "00:00:00 - 23:59:59"
                    }
                }
            }
        }
    },

    // 具体方法
    methods:{
        // 任务列表
        getConfigHash(ids){
            if (ids == undefined || ids == []){
                return []
            }
            let list = {}
            api.innerGet("/config/list", {ids: ids,page:1,size:ids.length}, (res)=>{
                if (!res.status){
                    this.$message.error(res.message);
                }else{
                    res.data.list.forEach(function (item) {
                        list[item.id] = item;
                    })
                }
            }, {async:false})
            return list
        },
        // 任务列表
        getList(){
            if (this.list.request){
                return this.$message.info('请求执行中,请稍等.');
            }
            this.list.request = true
            api.innerGet("/pipeline/list", this.listParam, (res)=>{
                this.list.request = false
                if (!res.status){
                    console.log("pipeline/list 错误", res)
                    return this.$message.error(res.message);
                }
                for (i in res.data.list){
                    let ratio = 0
                    if (res.data.list[i].top_number){
                        ratio = res.data.list[i].top_error_number / res.data.list[i].top_number
                    }
                    res.data.list[i].status = res.data.list[i].status.toString()
                    res.data.list[i].topRatio = 100 - ratio * 100
                }
                this.list.items = res.data.list;
                this.list.page = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.list.param.page = val
            this.getList()
        },
        handleClickTypeLabel(tab, event) {
            this.list.param.type = tab
            this.getList()
        },
        // 表单提交 流水线
        submitForm(){
            if (this.form.data.name == ''){
                return this.$message.error('请输入任务名称')
            }else if (this.form.data.spec == ''){
                return this.$message.error('请输入任务执行时间')
            }

            // 主要是强制类型
            let data = copyJSON(this.form.data)
            let body = {
                id: data.id,
                name: data.name,
                type: Number(data.type),
                spec: data.spec,
                config_ids: [],
                configs:data.configs,
                remark: data.remark,
                msg_set: data.msg_set,
                config_disable_action: Number(data.config_disable_action),
                config_err_action: Number(data.config_err_action),
            }
            data.configs.forEach(function (item,index) {
                body.config_ids.push(item.id)
                data.configs[index].status = Number(data.configs[index].status)
                if (data.configs[index].command.sql != null){
                    data.configs[index].command.sql.err_action = Number(data.configs[index].command.sql.err_action)
                }
            })

            api.innerPost("/pipeline/set", body, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.form.boxShow = false;
                this.form.data = this.initFormData()
                this.getList(); // 刷新当前页
            })
        },
        initFormData(){
            return  {
                id: 0,
                name: "",
                type: '2',
                config_ids:[], // 任务id集合
                configs:[], // 任务集合
                config_disable_action: '1',
                config_err_action: '1',
                msg_set: []
            }
        },
        // 编辑弹窗
        formBox(row){
            if (row == 0 || row == undefined){ // 添加显示
                this.form.boxShow = true
                this.form.boxTitle = '添加流水线'
                this.form.data = this.initFormData()
            }else if (row == -1){ // 关闭弹窗盒子
                this.form.boxShow = false
                return
            }else if (typeof row == 'object'){ // 编辑显示
                this.form.boxShow = true
                this.form.boxTitle = '编辑流水线'
                row = copyJSON(row)
                for (let i in row.msg_set){
                    row.msg_set[i] = this.msgSetBuildDesc(row.msg_set[i])
                }
                if (row){ // 拉取最新的任务信息进行展示
                    let configList = this.getConfigHash(row.config_ids)
                    for (let i=0; i<row.configs.length;i++){
                        let config = configList[row.configs[i].id]
                        if (config != undefined){
                            row.configs[i].name = config.name
                            row.configs[i].type_name = config.type_name
                            row.configs[i].status_name = config.status_name
                            row.configs[i].protocol_name = config.protocol_name
                        }
                    }
                }
                this.form.data = row
                if (this.form.data.msg_set == null){}
                this.form.data.status = this.form.data.status.toString()
                this.form.data.config_disable_action = this.form.data.config_disable_action.toString()
                this.form.data.config_err_action = this.form.data.config_err_action.toString()
                console.log("编辑：",this.form.data)
            }
            const that = this
            this.$nextTick(() => {
                this.sort = MySortable(document.getElementById("config-selected-box"), (oldIndex, newIndex)=>{
                    const oldlist = JSON.parse(JSON.stringify(this.form.data.configs));
                    const oldItem = oldlist.splice(oldIndex, 1)[0];
                    oldlist.splice(newIndex, 0, oldItem);

                    this.form.data.configs = []
                    that.$nextTick((t)=>{
                        this.form.data.configs= oldlist
                        console.log("拖拽后",this.form.data.configs, this, t)
                    })

                })
            })

        },
        // 改变状态
        changeStatus(row, newStatus){
            this.$confirm(newStatus==1? '确认关闭任务': '确认开启任务', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/pipeline/change_status", {id:row.id,status:Number(newStatus)}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    row.status = newStatus.toString()
                    row.status_name = newStatus == 1 ? '停用' :'激活'
                    return this.$message.success(res.message)
                })
            }).catch(()=>{
                // 取消操作
            })
        },

        getRegisterList(){
            this.register.boxShow = true;
            api.innerGet("/pipeline/register_list", {}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.register.items = res.data.list;
            })
        },
        // 成功率图标
        getTopIcon(total, ratio){
            // 没有执行
            if (total == 0){
                return "icon el-icon-remove-outline" // 0
            }
            // 成功率100
            if (ratio == 100){
                return "icon el-icon-circle-check success"
            }
            // 成功率66以上
            if (ratio < 100 && ratio > 66){
                return "icon el-icon-warning-outline warning"
            }
            // 成功率33以上
            if (ratio < 66 && ratio > 33){
                return "icon el-icon-circle-plus-outline danger"
            }
            // 成功率33一下
            return "icon el-icon-circle-close danger"
        },

        configLogBox(id, title){
            if (id == 0){
                this.log.show = false;
                this.log.id = 0;
            }else{
                this.log.id = id
                this.log.title = title+' 日志'
                this.log.show = true
            }
            console.log("configLogBox",id,title, this.log)
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
            this.msgSet.data = copyJSON(oldData)
        },
        msgSetDel(index){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('确认删除推送配置','提示',{
                type:'warning',
            }).then(()=>{
                this.form.msg_set.splice(index, 1);
            })

        },
        // 推送确认
        msgSetConfirm(){
            if (this.msgSet.data.msg_id <= 0){
                return this.$message.warning("请选择消息模板");
            }
            let data = this.msgSetBuildDesc(this.msgSet.data)

            if (this.msgSet.index < 0){
                this.form.data.msg_set.push(data)
            }else{
                this.form.data.msg_set[this.msgSet.index] = data
            }
            this.msgSet.show = false
            this.msgSet.index = -1
            this.msgSet.data = {}
        },
        // 构建消息设置描述
        msgSetBuildDesc(data){
            let item1 = this.msgSet.statusList.find(option => option.id === data.status);
            if (item1){
                data.status_name = item1.name
            }
            let descrition = '<i class="el-icon-bell"></i>当任务<span class="el-tag el-tag--small el-tag--light">'+item1.name+'</span>时'

            let item2 = this.dic.msg.find(option => option.id === data.msg_id)
            if (item2){
                data.msg_name = item2.name
                descrition += '，发送<span class="el-tag el-tag--small el-tag--light">'+item2.name+'</span>消息'
            }
            let item3 = this.dic.user.filter((option) => {
                return data.notify_user_ids.includes(option.id);
            }).map((item)=>{return item.name})
            if (item3.length > 0){
                data.notify_users_name = item3
                descrition += '，并且@人员<span class="el-tag el-tag--small el-tag--light">'+data.notify_users_name+'</span>'
            }
            data.descrition = descrition
            return data
        },
        statusClass(status){
            switch (Number(status)) {
                case 1:
                    return 'warning';
                case 2:
                    return 'primary';
                case 3:
                    return 'success';
                case 4:
                    return 'danger';
            }
        },
        // 枚举
        getDic(){
            api.dicList([Enum.dicUser, Enum.dicMsg],(res) =>{
                this.dic.user = res[Enum.dicUser]
                this.dic.msg = res[Enum.dicMsg]
            })
        },

        // 执行一下
        configRun(){
            let data = copyJSON(this.form.data)
            // 主要是强制类型
            let body = {
                id: data.id,
                name: data.name,
                type: Number(data.type),
                spec: data.spec,
                // protocol: Number(data.protocol),
                // command: data.command,
                remark: data.remark,
            }
            // body.command.sql.err_action = Number(body.command.sql.err_action)
            // body.command.sql.source.id = Number(body.command.sql.source.id)

            api.innerPost("/pipeline/run", body, (res)=>{
                if (!res.status){
                    return this.$message({
                        message: res.message,
                        type: 'error',
                        duration: 6000
                    })
                }
                return this.$message.success("ok."+res.data.result)
            })
        },
        // 任务盒子弹窗
        configSelectBox(e='show'){
            if (e == 'show'){ // 显示
                this.config.boxShow = true
            }else if (e == 'close'){ // 关闭
                this.config.boxShow = false
            }else if (e == 'confirm'){ // 提交
                console.log("选中元素",this.$refs.selection.selected)
                this.config.running = true
                this.form.data.configs.push(...this.$refs.selection.selected)
                this.config.boxShow = false
                this.config.running = false
            }
            console.log("任务盒子",this.config)
        },
        removeAt(idx) {
            this.$confirm('确定移除任务', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.form.data.configs.splice(idx, 1);
            }).catch(() => {/*取消*/});
        },
    }
})

Vue.component("MyPipeline", MyPipeline);