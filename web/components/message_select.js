var MyMessageSelect = Vue.extend({
    template: `<div class="message-select">
    <el-form :model="form" :inline="true" size="small">
        <el-form-item label="当执行">
            <el-select v-model="form.status" multiple style="width: 130px" placeholder="状态xx">
                <el-option v-for="(dic_v,dic_k) in msgSet.statusList" :label="dic_v.name" :value="dic_v.id"></el-option>
            </el-select>
            时
        </el-form-item>
        <el-form-item label="发送">
            <el-select v-model="form.msg_id" style="width: 180px" placeholder="xx模板">
                <el-option v-for="(dic_v,dic_k) in dic.msg" :label="dic_v.name" :value="dic_v.id"></el-option>
            </el-select>
            消息
        </el-form-item>
        <el-form-item label="并且@">
            <el-select v-model="form.notify_user_ids" multiple style="width: 210px" placeholder="人员">
                <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
            </el-select>
        </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
        <el-button @click="msgSet.show = false" size="small">取 消</el-button>
        <el-button type="primary" @click="msgSetConfirm()" size="small">确 定</el-button>
    </span>
</div>`,

    name: "MyMessageSelect",
    props: {
        value:Object
    },
    data(){
        return {
            form: {}, // 实际内容
            statusList:[{id:1,name:"错误"}, {id:2, name:"结束"}], //  {id:0,name:"开始"}
        }
    },
    // 模块初始化
    created(){
    },
    // 模块初始化
    mounted(){
        this.form = this.value
    },
    beforeDestroy(){},

    // 具体方法
    methods:{
        // 推送弹窗
        msgBoxShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (oldData == undefined || index < 0){
                oldData = {
                    status: [1],
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
        msgSetDel(index){
            if (this.request.disabled){
                return
            }
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
            if (this.form.msg_id <= 0){
                return this.$message.warning("请选择消息模板");
            }
            // 这里要补充name信息
            // let data = this.msgSetBuildDesc(this.form.data)

            // if (this.msgSet.index < 0){
            //     this.form.msg_set.push(data)
            // }else{
            //     this.form.msg_set[this.msgSet.index] = data
            // }
            // this.msgSet.show = false
            // this.msgSet.index = -1
            // this.msgSet.data = {}
            this.$emit('input', this.form); // 双向绑定，对应 props.value
            this.$emit('confirm',this.form)
        },
    }
})

Vue.component("MyMessageSelect", MyMessageSelect);