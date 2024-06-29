var MyConfigDetail = Vue.extend({
    template: `<el-main class="" style="max-width: 960px;margin: 15px auto 10px;">
        <div></div>
            <el-descriptions class="margin-top" :title="detail.name" :column="3" size="medium" border>
                <template slot="extra">
                  <el-button type="primary" size="small">操作</el-button>
                </template>
                
                <el-descriptions-item>
                  <template slot="label">状态</template>
                  <el-tooltip placement="top-start">
                    <div slot="content">{{detail.status_dt}}  {{detail.status_remark}}</div>
                    <span :class="statusClass(detail.status)">{{detail.status_name}}</span>
                  </el-tooltip>
                </el-descriptions-item>
                <el-descriptions-item>
                  <template slot="label">类型</template>
                  {{detail.protocol_name}}
                </el-descriptions-item>
                <el-descriptions-item>
                  <template slot="label">时间</template>
                  {{detail.type_name}}<el-divider direction="vertical"></el-divider>{{detail.spec}}
                </el-descriptions-item>
                
                <el-descriptions-item span="3">
                  <template slot="label">参数</template>
                  {{detail.var_fields}}
                </el-descriptions-item>
                
                <el-descriptions-item span="3">
                  <template slot="label">细节</template>
                  {{detail.command[detail.protocol_name]}}
                </el-descriptions-item>
                
                <el-descriptions-item>
                  <template slot="label">备注</template>
                  {{detail.remark}}
                </el-descriptions-item>
            </el-descriptions>
        <el-row>
            <h3>执行日志</h3>
            <my-config-log :tags="log.tags"></my-config-log>
        </el-row>
    </el-main>`,

    name: "MyConfigDetail",
    props: {
        request:{
            id:Number,
            type:String,
        }
    },
    data(){
        return {
            req:{id:0,type:''},
            detail:{},
            log:{
                tags:{}
            }
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        console.log('config_detail', this.$route, this.$router)

        if (this.$route.query.id){
            this.req.id = Number(this.$route.query.id)
        }
        if (this.$route.query.type){
            this.req.type = this.$route.query.type.toString()
        }
        this.log.tags = {ref_id: this.req.id, component:'config'}
        this.getDetail()
    },
    // 具体方法
    methods:{
        getDetail(){
            if (!this.req.id){
                return
            }
            api.innerGet("/config/detail", {id: this.req.id}, (res)=>{
                if (!res.status){
                    return this.$message({
                        message: res.message,
                        type: 'error',
                        duration: 6000
                    })
                }
                this.detail = res.data
            })
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
    }
})

Vue.component("MyConfigDetail", MyConfigDetail);