package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/changzee/multiavatar-go"
)

func main() {
	_ = os.MkdirAll("output", 0755)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/avatar", handleAvatar)

	addr := ":8080"
	log.Printf("Multiavatar demo server listening on %s\n", addr)
	log.Printf("Open: http://localhost%s/", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlIndex)
}

func handleAvatar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := strings.TrimSpace(q.Get("name"))
	if name == "" {
		http.Error(w, "missing required 'name' parameter", http.StatusBadRequest)
		return
	}

	var opts []multiavatar.Option

	// transparent => WithoutBackground
	if parseBool(q.Get("transparent")) {
		opts = append(opts, multiavatar.WithoutBackground())
	}

	// Global theme
	if t := strings.TrimSpace(q.Get("theme")); t != "" {
		opts = append(opts, multiavatar.WithTheme(t))
	}

	// Per-part theme: eyes:C,top:A
	for part, val := range parseKVComma(q.Get("partTheme")) {
		opts = append(opts, multiavatar.WithPartTheme(part, val))
	}

	// Allowed themes per part: top:A|C
	for part, list := range parseKVList(q.Get("allowedThemes")) {
		opts = append(opts, multiavatar.WithAllowedThemes(part, list))
	}

	// Force versions per part: eyes:11,top:07
	for part, val := range parseKVComma(q.Get("partVersion")) {
		opts = append(opts, multiavatar.WithPartVersion(part, val))
	}

	// Allowed versions per part: eyes:03|11,top:01|03|07
	for part, list := range parseKVList(q.Get("allowedVersions")) {
		opts = append(opts, multiavatar.WithAllowedVersions(part, list))
	}

	// Color overrides: env,clo,mouth,head,eyes,top with '|' separated values
	addColorOverrides(&opts, "env", q.Get("env"))
	addColorOverrides(&opts, "clo", q.Get("clo"))
	addColorOverrides(&opts, "mouth", q.Get("mouth"))
	if v := strings.TrimSpace(q.Get("head")); v != "" {
		opts = append(opts, multiavatar.WithSkinColor(v))
	}
	addColorOverrides(&opts, "eyes", q.Get("eyes"))
	addColorOverrides(&opts, "top", q.Get("top"))

	// Disable parts: top|eyes|clo|mouth|head|env
	for _, p := range splitList(q.Get("withoutPart")) {
		switch p {
		case "env", "clo", "head", "mouth", "eyes", "top":
			opts = append(opts, multiavatar.WithoutPart(p))
		}
	}

	svg := multiavatar.Generate(name, opts...)
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(svg))
}

// Helpers

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

// Parse "key:val,key2:val2" into map[key]val
func parseKVComma(s string) map[string]string {
	res := make(map[string]string)
	s = strings.TrimSpace(s)
	if s == "" {
		return res
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key != "" && val != "" {
			res[key] = val
		}
	}
	return res
}

// Parse "key:a|b|c,key2:x|y" into map[key][]string
func parseKVList(s string) map[string][]string {
	res := make(map[string][]string)
	s = strings.TrimSpace(s)
	if s == "" {
		return res
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		raw := strings.TrimSpace(kv[1])
		list := splitList(raw)
		if key != "" && len(list) > 0 {
			res[key] = list
		}
	}
	return res
}

func splitList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	items := strings.Split(s, "|")
	var out []string
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it != "" {
			out = append(out, it)
		}
	}
	return out
}

func addColorOverrides(opts *[]multiavatar.Option, part string, raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}
	colors := splitList(raw)
	if len(colors) == 0 {
		return
	}
	switch part {
	case "env":
		*opts = append(*opts, multiavatar.WithEnvColor(colors[0]))
	case "clo":
		*opts = append(*opts, multiavatar.WithClothesColors(colors...))
	case "mouth":
		*opts = append(*opts, multiavatar.WithMouthColors(colors...))
	case "eyes":
		*opts = append(*opts, multiavatar.WithEyesColors(colors...))
	case "top":
		*opts = append(*opts, multiavatar.WithTopColors(colors...))
	}
}

const htmlIndex = `<!doctype html>
<html lang="zh-cn">
<head>
  <meta charset="utf-8" />
  <title>multiavatar-go Demo</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <style>
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, sans-serif; margin: 0; padding: 24px; background: #f6f7f9; color: #222; }
    h1 { margin: 0 0 16px; font-size: 20px; }
    .container { display: grid; grid-template-columns: 360px 1fr; gap: 24px; align-items: start; }
    fieldset { border: 1px solid #ddd; border-radius: 8px; padding: 12px; background: #fff; }
    legend { padding: 0 6px; font-weight: 600; }
    label { display: block; margin: 8px 0 4px; font-size: 12px; color: #555; }
    input, select, button, textarea { width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 6px; font-size: 14px; }
    input[type="checkbox"] { width: auto; }
    .row { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
    .actions { display: flex; gap: 8px; margin-top: 12px; }
    button { cursor: pointer; background: #2563eb; color: #fff; border: none; }
    button.secondary { background: #6b7280; }
    #previewWrap { background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 16px; min-height: 260px; display:flex; align-items:center; justify-content:center; }
    #preview { width: 320px; height: 320px; }
    .hint { font-size: 12px; color: #666; }
    code { background: #eef; padding: 1px 4px; border-radius: 4px; }
  </style>
</head>
<body>
  <h1>multiavatar-go Demo</h1>
  <div class="container">
    <form id="form">
      <fieldset>
        <legend>基础</legend>
        <label>name（必填）</label>
        <input id="name" placeholder="例如：Alice" />
        <label>theme</label>
        <select id="theme">
          <option value="">默认</option>
          <option value="A">A</option>
          <option value="B">B</option>
          <option value="C">C</option>
        </select>
        <label class="row"><span><input type="checkbox" id="transparent" /> 透明背景</span></label>
      </fieldset>

      <fieldset>
        <legend>版本与主题筛选</legend>
        <label>partVersion（如 eyes:11,top:07）</label>
        <input id="partVersion" placeholder="eyes:11,top:07" />
        <label>allowedVersions（如 eyes:03|11,top:01|03|07）</label>
        <input id="allowedVersions" placeholder="eyes:03|11,top:01|03|07" />
        <label>partTheme（如 eyes:C,top:A）</label>
        <input id="partTheme" placeholder="eyes:C,top:A" />
        <label>allowedThemes（如 top:A|C,eyes:B|C）</label>
        <input id="allowedThemes" placeholder="top:A|C,eyes:B|C" />
        <label>withoutPart（禁用部件，例：top|eyes）</label>
        <input id="withoutPart" placeholder="top|eyes" />
      </fieldset>

      <fieldset>
        <legend>颜色</legend>
        <div class="row">
          <div>
            <label>env（背景色）</label>
            <input id="env" type="color" value="#eeeeee" />
          </div>
          <div>
            <label>head（肤色）</label>
            <input id="head" type="color" value="#f2c280" />
          </div>
        </div>
        <label>clo（衣服颜色，多个用 | 分隔，如 #333|#fff）</label>
        <input id="clo" placeholder="#333|#fff" />
        <label>mouth（嘴部颜色，多个用 | 分隔）</label>
        <input id="mouth" placeholder="#000|#fff" />
        <label>eyes（眼睛颜色，多个用 | 分隔）</label>
        <input id="eyes" placeholder="#000|#fff" />
        <label>top（发型颜色，多个用 | 分隔）</label>
        <input id="top" placeholder="#333|#ff0" />
        <div class="hint">提示：多颜色使用 | 分隔；页面会自动 URL 编码 #。</div>
      </fieldset>

      <div class="actions">
        <button type="button" onclick="preview()">预览</button>
        <button type="button" class="secondary" onclick="resetForm()">重置</button>
        <button type="button" class="secondary" onclick="download()">下载SVG</button>
      </div>
      <div class="hint" id="urlHint"></div>
    </form>

    <div id="previewWrap">
      <div id="preview"></div>
    </div>
  </div>

  <script>
    function buildURL() {
      const params = new URLSearchParams();
      const name = document.getElementById('name').value.trim();
      if (!name) { alert('请填写 name'); return null; }
      params.set('name', name);

      const theme = document.getElementById('theme').value.trim();
      if (theme) params.set('theme', theme);

      const transparent = document.getElementById('transparent').checked;
      if (transparent) params.set('transparent', 'true');

      const ids = ['partVersion','allowedVersions','partTheme','allowedThemes','withoutPart','env','head','clo','mouth','eyes','top'];
      for (const id of ids) {
        const v = document.getElementById(id).value.trim();
        if (v) params.set(id, v);
      }

      const url = '/avatar?' + params.toString();
      document.getElementById('urlHint').innerHTML = '请求：<code>' + url.replace(/&/g,'&') + '</code>';
      return url;
    }

    async function preview() {
      const url = buildURL();
      if (!url) return;
      try {
        const resp = await fetch(url);
        if (!resp.ok) throw new Error('请求失败: ' + resp.status);
        const svg = await resp.text();
        document.getElementById('preview').innerHTML = svg;
      } catch (e) {
        alert(e.message);
      }
    }

    function resetForm() {
      document.getElementById('form').reset();
      document.getElementById('preview').innerHTML = '';
      document.getElementById('urlHint').textContent = '';
    }

    function download() {
      const url = buildURL();
      if (!url) return;
      const a = document.createElement('a');
      a.href = url;
      a.download = 'avatar.svg';
      document.body.appendChild(a);
      a.click();
      a.remove();
    }
  </script>
</body>
</html>`
