<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <title>个人客户列表</title>
    <link rel="stylesheet" href="/static/plug/layui/css/layui.css">
</head>
<body>
<script type="text/javascript" src="/static/plug/layui/layui.js"></script>
<form class="layui-form" action="/admin/save"  method="post" style="margin:20px">
    <div class="layui-form-item">
        <label class="layui-form-label">类型</label>
        <div class="layui-input-block">
            <input type="radio" name="types" value="1" title="博文"  {{if .post.Types}} checked {{end}}>
            <input type="radio" name="types" value="0" title="下载"  {{if .post.Types}} {{else}} checked {{end}}>
        </div>
    </div>
    <input type="text" name="Id" style="visibility: hidden;" value = "{{.post.Id}}">

    <div class="layui-form-item">
        <label class="layui-form-label">标题：</label>
        <div class="layui-input-block">
            <input type="text" name="title" required value="{{.post.Title}}" lay-verify="required" placeholder="请输入标题" autocomplete="off" class="layui-input">
        </div>
    </div>

    <div class="layui-form-item">
        <label class="layui-form-label">类别</label>
        <div class="layui-input-block">
            <select name="cate_id" lay-verify="required">
                {{if .categorys}}
                        {{range .categorys}}
                            <option value="{{.Id}}" {{if eq $.post.CategoryId .Id}} selected {{end}}>{{.Name}}</option>
                        {{end}}
                    {{else}}
                        <option value=""></option>
                    {{end}}
            </select>
        </div>
    </div>

    <div class="layui-form-item">
        <label class="layui-form-label">加入首页</label>
        <div class="layui-input-block">
            <input type="checkbox" name="is_top" {{if .post.IsTop}} checked {{end}} value="1" title="置首" >
        </div>
    </div>

    <div class="layui-form-item">
        <label class="layui-form-label">链接</label>
        <div class="layui-input-block">
            <input type="text" name="url" lay-verify="url" value="{{.post.Url}}" placeholder="请输入下载链接" autocomplete="off" class="layui-input">
        </div>
    </div>


    <div class="layui-form-item">
        <label class="layui-form-label">标签</label>
        <div class="layui-input-block">
            <input type="text" name="tags" value="{{.post.Tags}}"  placeholder="标签，隔开" autocomplete="off" class="layui-input">
        </div>
    </div>


    <div class="layui-form-item">
        <label class="layui-form-label">标签</label>
        <div class="layui-input-block">
            <textarea name="info" placeholder="请输入内容" class="layui-textarea">{{.post.Info}}</textarea>
        </div>
    </div>



    <div class="layui-form-item">
        <label class="layui-form-label">图片</label>
        <div class="layui-input-block">
            <input name="Image" id="Image" value="{{.post.Image}}" placeholder="请输入内容" style="width: 300px; float: left"  class="layui-input">
            <input type="file" name="uploadname" style="float: left" lay-ext="jpg|png|gif" class="layui-upload-file">
        </div>
    </div>

    <div class="layui-form-item layui-form-text">
        <label class="layui-form-label">内容</label>
        <div class="layui-input-block">
            <script type="text/javascript" charset="utf-8">
                window.UEDITOR_HOME_URL = "/static/ueditor/";
            </script>
            <script type="text/javascript" src="/static/ueditor/ueditor.config.js"></script>
            <script type="text/javascript" src="/static/ueditor/ueditor.all.min.js"></script>
            <textarea id="content" name="content">{{.post.Content}}</textarea>
            <script type="text/javascript" charset="utf-8">
                var options = {"fileUrl":"/admin/article/upload",
                    "filePath":"","imageUrl":"/admin/article/upload","imagePath":"",
                    "initialFrameWidth":"90%",
                    "initialFrameHeight":"400",
                };
                var ue = UE.getEditor("content",options);
            </script>
        </div>
    </div>
    <div class="layui-form-item">
        <div class="layui-input-block">
            <button class="layui-btn" lay-submit lay-filter="formDemo">提交</button>
            <button type="reset" class="layui-btn layui-btn-primary">重置</button>
        </div>
    </div>
</form>
<script>
    //Demo
    layui.use('form', function(){


    });

    layui.use('upload', function(){
        layui.upload({
            url: '/admin/upload'
            ,success: function(res, input){
                if(res.code ==0 ){
                   document.getElementById("Image").value = res.message;
                }else{
                    layui.msg(res.message)
                }
            }
        });
    });

</script>

</body>
</html>