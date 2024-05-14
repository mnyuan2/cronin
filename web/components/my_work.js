var MyWork = Vue.extend({
    template: `<el-main class="my-work">
<!--
        先规划一下页面布局（整体仿造tapd，还有我们自己的手机端），
        # 各个环境名称（如果没有对应任务就不在哪上）（默认只有第一个标签页打开）
            流水线和任务混合在一起（sql写法要研究一下）
                点击跳转到详情页（新的标签页）
-->
        <el-row v-for="(group_v,group_k) in group_list">
            <el-row class="header">
                <i class="el-icon-caret-right my-table-icon"></i>
                {{group_v.env_title}} <el-divider direction="vertical"></el-divider> 
                {{group_v.join_type=='config'? '任务' : '流水线'}} <span style="color: #72767b">（{{group_v.total}}）</span>
            </el-row>
            <el-row class="body">body</el-row>
        </el-row>
    </el-main>`,

    name: "MyWork",
    props: {
        data_id:Number
    },
    data(){
        return {
            group_list: [],
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        this.getTables()
    },
    // 具体方法
    methods:{
        getTables(){
            api.innerGet('/work/table',{}, res =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.group_list = res.data.list
            })
        }
    }
})

Vue.component("MyWork", MyWork);