import{d as $,o as w,a,c as l,b as e,F as E,r as P,g as u,h as p,u as x,t as n,n as I}from"./index-DiD6N43a.js";import{_ as j}from"./_plugin-vue_export-helper-DlAUqK2U.js";const B={class:"api-view"},S={class:"resource-grid"},z={class:"resource-header"},D={class:"provider"},H=["title"],M={class:"description"},N={class:"quota"},T={class:"value"},q={class:"endpoint"},L={class:"resource-actions"},O=["href"],R=["onClick"],V=["onClick"],F={key:0,class:"code-examples"},K={class:"code-block"},Y={class:"code-header"},X=["onClick"],G={class:"code-block"},J={class:"code-header"},Q=["onClick"],U={class:"updated"},W={key:0,class:"empty"},Z=$({__name:"ResourcesView",setup(ee){const f=x(),c=p([]),_=p(!0),d=p({}),r=p({});async function k(){_.value=!0;try{const t=await(await fetch("/api/resources")).json();c.value=t.resources||[],y()}catch(o){console.error("Failed to load API resources:",o)}finally{_.value=!1}}async function y(){for(const o of c.value)try{const s=await(await fetch(`/api/resources/health?id=${o.id}`)).json();d.value[o.id]=s.status||"unknown"}catch{d.value[o.id]="unknown"}}function C(o){f.push(`/chat?model=${o}`)}function b(o){r.value[o]=!r.value[o]}function h(o){const t=o.model||"MODEL_NAME",s=o.provider==="阿里云"?"Authorization: Bearer $DASHSCOPE_API_KEY":"Authorization: Bearer $API_KEY";return`curl -X POST "${o.endpoint}" \\
  -H "Content-Type: application/json" \\
  -H "${s}" \\
  -d '{
    "model": "${t}",
    "messages": [{"role":"user","content":"你好"}]
  }'`}function v(o){const t=o.model||"MODEL_NAME";return`import requests

response = requests.post(
    "${o.endpoint}",
    headers={
        "Content-Type": "application/json",
        "Authorization": "Bearer $API_KEY"
    },
    json={
        "model": "${t}",
        "messages": [{"role": "user", "content": "你好"}]
    }
)
print(response.json())`}function m(o,t){const s=c.value.find(A=>A.id===o);if(!s)return;const i=t==="curl"?h(s):v(s);navigator.clipboard.writeText(i)}function g(o){return o?new Date(o).toLocaleString("zh-CN"):""}return w(k),(o,t)=>(a(),l("div",B,[t[4]||(t[4]=e("div",{class:"api-header"},[e("h2",null,"免费 API 资源清单"),e("p",{class:"subtitle"},"当前可用的免费大模型 API 资源")],-1)),e("div",S,[(a(!0),l(E,null,P(c.value,s=>(a(),l("div",{key:s.id,class:"resource-card"},[e("div",z,[e("div",null,[e("h3",null,n(s.name),1),e("span",D,n(s.provider),1)]),e("span",{class:I(["status-dot",d.value[s.id]||"unknown"]),title:d.value[s.id]},null,10,H)]),e("p",M,n(s.description),1),e("div",N,[t[0]||(t[0]=e("span",{class:"label"},"免费额度：",-1)),e("span",T,n(s.free_quota),1)]),e("div",q,[t[1]||(t[1]=e("span",{class:"label"},"端点：",-1)),e("code",null,n(s.endpoint),1)]),e("div",L,[s.doc_url?(a(),l("a",{key:0,href:s.doc_url,target:"_blank",class:"doc-link"},"查看文档",8,O)):u("",!0),s.model?(a(),l("button",{key:1,class:"trial-btn",onClick:i=>C(s.model)}," 一键试用 ",8,R)):u("",!0),e("button",{class:"code-toggle-btn",onClick:i=>b(s.id)},n(r.value[s.id]?"收起示例":"代码示例"),9,V)]),r.value[s.id]?(a(),l("div",F,[e("div",K,[e("div",Y,[t[2]||(t[2]=e("span",null,"curl",-1)),e("button",{class:"copy-btn",onClick:i=>m(s.id,"curl")},"复制",8,X)]),e("pre",null,[e("code",null,n(h(s)),1)])]),e("div",G,[e("div",J,[t[3]||(t[3]=e("span",null,"Python",-1)),e("button",{class:"copy-btn",onClick:i=>m(s.id,"python")},"复制",8,Q)]),e("pre",null,[e("code",null,n(v(s)),1)])])])):u("",!0),e("div",U,"上次更新："+n(g(s.last_updated)),1)]))),128)),!c.value.length&&!_.value?(a(),l("div",W,"暂无 API 资源")):u("",!0)])]))}}),te=j(Z,[["__scopeId","data-v-f92a2cb3"]]);export{te as default};
