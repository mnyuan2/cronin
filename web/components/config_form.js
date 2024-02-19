var MyConfigForm = Vue.extend({
    template: `<div class="config-form">
    
</div>`,
    name: "MyConfigForm",
    props: {
        config_id:Number,
    },
    data(){
        return {
            config_id:0,
            form:{},
            hintSpec: "* * * * * *",
            sqlSourceBoxShow: false,
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
        }
    },
    // 模块初始化
    created(){},
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
                    },
                },
                msg_set: []
            }
        },
    }
})

Vue.component("MyConfigForm", MyConfigForm);