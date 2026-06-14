#!/bin/bash
sudo mysql -e "
CREATE DATABASE IF NOT EXISTS ai_watcher CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ai_watcher;
SOURCE /projects/mimo/ai-watcher/migrations/001_init.sql;
INSERT INTO api_resources (name, provider, description, endpoint, free_quota, doc_url) VALUES
('Qwen-Turbo','阿里云','轻量模型','https://dashscope.aliyuncs.com/api/v1','100万tokens/月','https://help.aliyun.com/zh/dashscope'),
('DeepSeek V3','DeepSeek','大语言模型','https://api.deepseek.com/v1','500万tokens/月','https://platform.deepseek.com'),
('GLM-4-Flash','智谱AI','GLM-4轻量','https://open.bigmodel.cn/api/paas/v4','100万tokens/月','https://bigmodel.cn'),
('ERNIE-Speed','百度','文心一言轻量版','https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k','每天10000次','https://cloud.baidu.com/doc/WENXINWORKSHOP'),
('Spark Lite','科大讯飞','星火轻量版','https://spark-api-open.xf-yun.com/v1/chat/completions','每天5000次','https://www.xfyun.cn');
" && echo "数据库初始化完成" && cd /projects/mimo/ai-watcher && nohup ./ai-watcher > /tmp/ai-watcher.log 2>&1 & sleep 2 && curl -s -o /dev/null -w "HTTP %{http_code}" http://localhost:8080 && echo " 启动成功! 访问 http://localhost:8080" || (echo " 启动失败" && cat /tmp/ai-watcher.log)
