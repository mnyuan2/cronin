var MyTrace = Vue.extend({
    template: `<div class="trace-page">
    <div class="median" style="left: 300px;"></div>
    <el-row class="trace-header">
        <div class="label">操作</div>
        <div class="total-desc" v-html="total_desc"></div>
    </el-row>
    <el-table :show-header="false" :data="traces" row-key="span_id" default-expand-all :tree-props="{children:'children',hasChildren:'hasChildren'}" class="trace-wrapper">
        <template #empty>空</template>
        <el-table-column width="300">
            <template slot-scope="scope">
              <div :class="(scope.row.status==1?'is-light':'')+ ' span-label'" @click="showSpan(scope.row.bar, 'detail_show')">
                <i class="el-alert__icon el-icon-error" v-if="scope.row.status==1"></i>
                <span>{{scope.row.operation}}</span>
                <i class="el-alert__icon el-icon-loading info-2" v-if="!scope.row.bar.finish" title="执行中"></i>
              </div>
            </template>
        </el-table-column>
        <el-table-column>
            <template slot-scope="scope">
                <div class="span-bar-wrapper" @click="showSpan(scope.row.bar, 'detail_show')">
                    <div class="span-bar" :style="scope.row.bar.style">
                        &nbsp;
                        <div class="span-bar-label" :style=" scope.row.bar.width>80? 'right: 0;' : scope.row.bar.left>40 ? 'right: 100%;' : 'left: 100%;'">{{scope.row.bar.label}}</div>
                    </div>
                </div>
                <div class="span-detail" v-if="scope.row.bar.detail_show">
                
                    <!-- tags -->
                    <div class="span-wrapper">
                        <div class="span-header" @click="showSpan(scope.row.bar, 'tags_show')">
                            <i :class="scope.row.bar.tags_show? 'el-icon-arrow-down':'el-icon-arrow-right'"></i>
                            <strong>Tags: </strong>
                            <el-breadcrumb v-if="!scope.row.bar.tags_show">
                                <el-breadcrumb-item v-for="(tag_v,tag_i) in scope.row.tags">{{tag_v.key}}={{tag_v.value.value}}</el-breadcrumb-item>
                            </el-breadcrumb>
                        </div>
                        <el-table :data="scope.row.tags" stripe :show-header="false" style="width: 100%" v-if="scope.row.bar.tags_show">
                            <el-table-column prop="key" width="180"> </el-table-column>
                            <el-table-column prop="value.value"> </el-table-column>
                        </el-table>
                    </div>
                  
                    <!-- logs -->
                    <div class="span-logs-box">
                        <div class="span-header" @click="showSpan(scope.row.bar, 'logs_show')">
                            <i :class="scope.row.bar.logs_show? 'el-icon-arrow-down':'el-icon-arrow-right'"></i>
                            <strong>Logs: </strong>
                            <el-breadcrumb>
                                <el-breadcrumb-item>{{scope.row.logs.length}}</el-breadcrumb-item>
                            </el-breadcrumb>
                        </div>
                        
                        <div class="span-wrapper"  v-if="scope.row.bar.logs_show" v-for="(log_v, log_i) in scope.row.logs">
                            <div class="span-header" @click="showSpan(scope.row.bar.log[log_i], 'show')">
                                <i :class="scope.row.bar.log[log_i].show? 'el-icon-arrow-down':'el-icon-arrow-right'"></i>
                                <strong>{{durationTransform(log_v.timestamp-scope.row.timestamp)}}: </strong>
                                <el-breadcrumb v-if="!scope.row.bar.log[log_i].show">
                                    <el-breadcrumb-item v-for="(log_f_v,log_f_i) in log_v.attributes">{{log_f_v.key}}={{log_f_v.value.value}}</el-breadcrumb-item>
                                </el-breadcrumb>
                            </div>
                            <el-table :data="log_v.attributes" stripe :show-header="false" style="width: 100%" v-if="scope.row.bar.log[log_i].show">
                                <el-table-column prop="key" width="180"> </el-table-column>
                                <el-table-column>
                                    <template slot-scope="scope">
                                        <!-- json格式化 -->
                                        <pre v-if="scope.row.value.type=='STRING' && isJSON(scope.row.value.value)" style="max-height: 350px;overflow: auto;line-height: 100%;">{{JSON.parse(scope.row.value.value)}}</pre>
                                        <!-- 异常文本 换行格式化 v-else-if="scope.row.key=='error.panic' || scope.row.key=='sql' || scope.row.key=='statement'" -->
                                        <pre v-else style="max-height: 350px;overflow: auto;line-height: 120%;">{{scope.row.value.value}}</pre>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </div>
                    </div>
                </div>
            </template>
      </el-table-column>
    </el-table>
    
               
</div>`,
    name: "MyConfigLog",
    props: {
        trace_id:String,// 与 entry_id 互斥
        job:{
            ref_id: Number,
            entry_id: Number,
            trace_id: String,
        },// 与 trace_id 互斥
    },
    data(){
        return {
            trace_id:"",
            job:{},
            traces:[], // 踪迹
            total_desc: "", // 合计描述
            job_task:0,
            last_span_map:{}
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){},
    // 实例销毁前
    beforeDestroy(){
        this.removeJob()
    },
    watch:{
        trace_id:{
            immediate: true, // 解决首次赋值不触发的情况
            handler: function (newVal,oldVal){
                console.log("trace trace_id",newVal, oldVal)
                if (newVal){
                    this.getTrace("/log/traces", {trace_id: newVal})
                }
            },
        },
        job:{
            immediate: true,
            handler: function (newVal,oldVal){
                console.log("trace job",newVal, oldVal)
                if (newVal && !this.trace_id){
                    this.removeJob()
                    this.job_task = setInterval(()=>{this.getTrace("/job/traces", newVal)}, 4000)
                }
            },
        }
    },

    // 具体方法
    methods:{
        // 踪迹展示
        getTrace(path, param){
            api.innerGet(path, param, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                const micro = Date.now()*1000;

                let list = res.data.list[0].Spans
                let lastSpan = list.slice(-1)[0]
                let startSpan = list[0]
                // 截止的截止-开始时间 = 总耗时
                let totalDuration = lastSpan.timestamp+(lastSpan.duration>=0?lastSpan.duration: micro - startSpan.timestamp)-startSpan.timestamp
                let startDuration = list[0].duration>=0?list[0].duration: micro - startSpan.timestamp
                if (totalDuration<startDuration){
                    totalDuration = startDuration // 这里应该是等于最大的一个耗时
                }
                if (startSpan.duration>=0){
                    this.removeJob()
                }
                let span_map = {}
                for (let span of list){
                    // 视图信息控制参数
                    span.bar = {
                        left: (span.timestamp-startSpan.timestamp)/totalDuration*100, // 节点起始 = 当前节点开始-开始的开始
                        detail_show: false,
                        tags_show: false,
                        logs_show: false,
                        finish: true,
                        log:[], // 单个视图元素
                    }
                    if (span.duration<0){
                        span.duration = micro - span.timestamp
                        span.bar.finish = false
                    }
                    span.bar.width = span.duration/totalDuration*100 // 节点截止 = 当前耗时占比总耗时的比例·宽
                    span.bar.label = durationTransform(span.duration) // 这里后面要给个单位出来（前端转）
                    span.logs.forEach(function (item) {
                        span.bar.log.push({show:false})
                    })
                    span.bar.style = 'background: #37be5f; left: '+(span.bar.left> 99.9? 99.9 : span.bar.left)+'%; width: '+(span.bar.width<0.1? '1px': span.bar.width+'%')+';'
                    span_map[span.span_id] = span
                    if (this.last_span_map[span.span_id]){ // 重复请求时，继承原UI属性
                        let last_bar = this.last_span_map[span.span_id].bar
                        span.bar.detail_show = last_bar.detail_show
                        span.bar.tags_show = last_bar.tags_show
                        span.bar.logs_show = last_bar.logs_show
                        last_bar.log.forEach(function (item, index) {
                            span.bar.log[index] = item
                        })
                    }
                }
                this.last_span_map = span_map

                let trace = arrayToTree(list, 'span_id', 'parent_span_id', 'children')
                console.log(trace, totalDuration)
                this.traces = trace;
                this.total_desc = "<span>开始时间:<b>"+getDatetimeString(new Date(startSpan.timestamp/1000))+
                    "</b></span>  <span>耗时:<b>"+ durationTransform(totalDuration) +
                    "</b></span>  <span>节点数:<b>"+list.length +"</b><span>"
            })
        },
        // 节点显示ui控制
        showSpan(row, field){
            if (typeof field !== 'string' || field == ""){
                alert('请指定字段')
                retrun
            }
            if(row[field]==false){
                row[field] = true
            }else{
                row[field] = false
            }
        },
        // log格式化
        formatLog(log){
            if (log.value.type=='STRING' && isJSON(log.value.value)){
                return '<pre>'+log.value.value+'</pre>'
            }else{
                return log.value.value
            }
        },
        // 移除任务
        removeJob(){
            if (this.job_task> 0){
                clearInterval(this.job_task)
                this.job_task = 0
            }
        }
    }
})

Vue.component("MyTrace", MyTrace);