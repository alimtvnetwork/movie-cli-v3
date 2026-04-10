# Modal Structure & Components

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

## Component Hierarchy

```
GlobalErrorModal.tsx
├── DialogHeader (error code, timestamp, queue navigation)
├── Section Toggle (Backend / Frontend buttons)
├── ScrollArea
│   ├── BackendSection.tsx (when activeSection === "backend")
│   │   ├── Tabs
│   │   │   ├── Overview (error message, request, timing, badges)
│   │   │   ├── Log (error.log.txt content)
│   │   │   ├── Execution (Go call chain table + session logs)
│   │   │   ├── Stack (Go + PHP stack traces, session diagnostics)
│   │   │   ├── Session (SessionLogsTab — logs, request, response, stack trace)
│   │   │   ├── Request (RequestDetails — request chain visualization)
│   │   │   └── Traversal (TraversalDetails — endpoint flow + methods stack)
│   │   └── (Internal sub-components: OverviewContent, ErrorLogContent, etc.)
│   │
│   └── FrontendSection.tsx (when activeSection === "frontend")
│       ├── Tabs
│       │   ├── Overview (trigger context, message, call chain, click path)
│       │   ├── Stack (parsed/raw JS stack, React execution chain)
│       │   ├── Context (full error context JSON)
│       │   └── Fixes (suggested fixes by error code)
│       └── (Internal sub-components)
│
├── DialogFooter
│   ├── DownloadDropdown (ZIP bundle, error.log.txt, log.txt, report.md)
│   ├── Close button
│   └── CopyDropdown (Split Button: main click = compact report, chevron dropdown = full report, with backend logs, error.log.txt, log.txt)
│
└── ErrorModalTypes.ts (shared types: PHPStackFrame, AppInfo, SectionCommonProps)
```

---

## Visual Layout Diagrams

### Full Modal Layout (Desktop: 95vw × 95vh)

```
┌─────────────────────────────────────────────────────────────────────┐
│ ┌─ DialogHeader ──────────────────────────────────────────────────┐ │
│ │ [E5001]  Failed to enable plugin   2026-02-09 14:32:01         │ │
│ │                                          [◀ 1/3 ▶] [Copy All] │ │
│ └─────────────────────────────────────────────────────────────────┘ │
│ ┌─ Section Toggle ────────────────────────────────────────────────┐ │
│ │  [ ● Backend ]  [ ○ Frontend ]                                 │ │
│ └─────────────────────────────────────────────────────────────────┘ │
│ ┌─ ScrollArea (flex-1) ──────────────────────────────────────────┐ │
│ │ ┌─ Tab Bar ──────────────────────────────────────────────────┐ │ │
│ │ │ Overview │ Log │ Execution │ Stack │ Session │ Request │ Traversal │
│ │ └────────────────────────────────────────────────────────────┘ │ │
│ │                                                                │ │
│ │  ┌─ Active Tab Content ─────────────────────────────────────┐  │ │
│ │  │                                                          │  │ │
│ │  │  (Tab-specific content rendered here)                    │  │ │
│ │  │                                                          │  │ │
│ │  └──────────────────────────────────────────────────────────┘  │ │
│ └────────────────────────────────────────────────────────────────┘ │
│ ┌─ DialogFooter ─────────────────────────────────────────────────┐ │
│ │  [▼ Download]                              [Close] [▼ Copy]   │ │
│ └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Backend Section — Overview Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Error Banner (red) ────────────────────────────────────────┐ │
│  │ ⚠ Backend Error: Failed to fetch plugin details from site  │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Delegated Info Banner (blue, NEW v2.0.0) ──────────────────┐ │
│  │ ℹ Endpoint 'snapshots' is not enabled in plugin settings.   │ │
│  │   (from DelegatedRequestServer.AdditionalMessages)          │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Request Info ──────────────────────────────────────────────┐ │
│  │  Method: POST   Endpoint: /api/v1/plugins/enable            │ │
│  │  Status: 500    Site: https://example.com                   │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌─ Timing ───────────────────────────────────────────────────┐  │
│  │  Requested At:           /api/v1/plugins/enable             │  │
│  │  Delegated At:           https://site.com/wp-json/...       │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Availability Badges ──────────────────────────────────────┐  │
│  │  [✓ Session] [✓ Stack Traces] [✓ Delegated Info] [✓ Exec] │  │
│  └─────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Backend Section — Stack Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Go Backend Stack (blue-themed) ───────────────────────────┐  │
│  │  site_handlers.go:327  handlers.EnableRemotePlugin          │  │
│  │  service.go:1245       site.(*Service).EnableRemotePlugin   │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Server Stack (purple-themed, NEW v2.0.0) ───────┐  │
│  │  ┌─ Header ─────────────────────────────────────────────┐   │  │
│  │  │ 🟣 Delegated Server  GET  403                        │   │  │
│  │  │ https://example.com/wp-json/riseup.../snapshots/...  │   │  │
│  │  └──────────────────────────────────────────────────────┘   │  │
│  │  Stack Trace:                                               │  │
│  │  #0 riseup-asia-uploader.php(1098): Logger->error()         │  │
│  │  #1 class-wp-hook.php(341): Plugin->enrichError()          │  │
│  │  #2 plugin.php(205): WP_Hook->apply_filters()               │  │
│  │  Response:                                                  │  │
│  │  ▸ { "code": "rest_forbidden", "message": "...", ... }     │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ PHP Delegated Stack (orange-themed, legacy) ───────────────┐  │
│  │  PHP Fatal error: Class 'PDO' not found in plugin-mgr.php  │  │
│  │  #0 endpoints.php(15): PluginManager->connect()             │  │
│  │  #1 {main}                                                  │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ PHP Structured Frames (table) ────────────────────────────┐  │
│  │  #  │ Function                    │ File              │ Line│  │
│  │  0  │ PluginManager::connect()    │ plugin-mgr.php    │ 42  │  │
│  │  1  │ handle_enable()             │ endpoints.php     │ 15  │  │
│  └─────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Session Diagnostics (auto-fetched) ───────────────────────┐  │
│  │  Go frames: 3 │ PHP frames: 2 │ stacktrace.txt: available  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Backend Section — Request Tab (Chain Visualization)

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Node 1: React → Go ──────────────────────────────────────┐  │
│  │  🔵 [React → Go]  [POST]  [500]                           │  │
│  │  /api/v1/plugins/enable                                    │  │
│  │  ▸ Request Body: { "slug": "my-plugin", "SiteId": 1 }     │  │
│  └──────────┬─────────────────────────────────────────────────┘  │
│             │ (vertical connector line)                           │
│  ┌──────────┴─────────────────────────────────────────────────┐  │
│  │  🟠 [Go → Delegated]  [GET]  [403]                        │  │
│  │  https://example.com/wp-json/riseup.../v1/snapshots/...    │  │
│  │  ▸ Request Body: (none — GET)                              │  │
│  └──────────┬─────────────────────────────────────────────────┘  │
│             │ (vertical connector line)                           │
│  ┌──────────┴─────────────────────────────────────────────────┐  │
│  │  🟣 [Delegated Response]  (NEW v2.0.0)                     │  │
│  │  ▸ Response: { "code": "rest_forbidden", "message": ... }  │  │
│  │  ▸ Stack Trace: #0 riseup-asia-uploader.php(1098)...       │  │
│  │  ▸ Additional: Endpoint 'snapshots' is not enabled...      │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Environment ──────────────────────────────────────────────┐  │
│  │  API Base: http://localhost:8080   VITE_API_URL: ...        │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Backend Section — Traversal Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Endpoint Flow (3-hop, NEW v2.0.0) ─────────────────────────┐ │
│  │  [React] http://localhost:8080                               │ │
│  │    ──▸                                                      │ │
│  │  [Go] /api/v1/sites/1/snapshots/settings                    │ │
│  │    ──▸                                                      │ │
│  │  [Delegated] https://site.com/wp-json/riseup.../settings    │ │
│  │              GET → 403                                      │ │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Methods Stack (table) ────────────────────────────────────┐  │
│  │  #  │ Method                          │ File            │ Ln│  │
│  │  1  │ handlers.handleSiteActionById   │ handler_factory │107│  │
│  │  2  │ api.SessionLogging              │ session_log     │107│  │
│  │  3  │ api.Recovery                    │ middleware      │245│  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Server Details (purple, NEW v2.0.0) ─────────────┐ │
│  │  Endpoint: https://site.com/wp-json/riseup.../settings      │ │
│  │  Method: GET │ Status: 403                                  │ │
│  │  Stack Trace:                                                │ │
│  │    #0 riseup-asia-uploader.php(1098): Logger->error()       │ │
│  │    #1 class-wp-hook.php(341): enrichErrorResponse()       │ │
│  │  Additional: Endpoint not enabled in plugin settings        │ │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Delegated Service Error Stack (orange, legacy) ────────────┐ │
│  │  PHP Fatal error: Class 'PDO' not found...                  │ │
│  │  #0 endpoints.php(15): PluginManager->connect()             │ │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Frontend Section — Overview Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  ┌─ Trigger Context ──────────────────────────────────────────┐  │
│  │  Component: PluginCard  →  Action: enable_clicked           │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Message ──────────────────────────────────────────────────┐  │
│  │  Failed to enable plugin "my-plugin" on site                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ Call Chain ───────────────────────────────────────────────┐  │
│  │  PluginsPage                                                │  │
│  │    └─ usePluginActions.enable                               │  │
│  │        └─ api.post("/api/v1/plugins/enable")                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ User Interaction Path (last 10 clicks) ───────────────────┐  │
│  │  14:31:55  PluginsPage     "Plugins" tab       click  /    │  │
│  │  14:31:58  PluginCard      "Enable" button     click  /    │  │
│  │  14:32:01  PluginCard      "Confirm" button    click  /    │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Frontend Section — Stack Tab

```
┌──────────────────────────────────────────────────────────────────┐
│  [● Parsed] [○ Raw]                    [□ Show internal frames] │
│                                                                  │
│  ┌─ Parsed Stack Frames (table) ──────────────────────────────┐  │
│  │  #  │ Function              │ File                │ Line   │  │
│  │  0  │ enablePlugin          │ usePluginActions.ts  │ 45    │  │
│  │  1  │ handleClick           │ PluginCard.tsx       │ 112   │  │
│  │  2  │ callCallback          │ react-dom.js         │ 3942  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌─ React Execution Chain ────────────────────────────────────┐  │
│  │  [render] PluginsPage                           14:31:50   │  │
│  │  [effect] usePluginActions                      14:31:51   │  │
│  │  [handler] enablePlugin                         14:32:01   │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  Error Location: usePluginActions.ts:45 in enablePlugin()        │
└──────────────────────────────────────────────────────────────────┘
```

### DialogFooter — Action Menus

The **Copy** button uses a **Split Button** pattern: the main button area copies the **Compact Report** instantly (no API call), while the chevron arrow opens a dropdown with all copy options.

```
┌──────────────────────────────────────────────────────────────────┐
│  [▼ Download]                       [Close]  [ Copy ][▼]        │
│  ┌─────────────────┐                    ┌──────────────────────┐│
│  │ Full Bundle (ZIP)│        main click: │ (copies compact)     ││
│  │ error.log.txt    │        chevron ▼:  │ Compact Report       ││
│  │ log.txt          │                    │ Full Report          ││
│  │ Report (.md)     │                    │ With Backend Logs    ││
│  └─────────────────┘                    │ error.log.txt        ││
│                                         │ log.txt              ││
│                                         └──────────────────────┘│
└──────────────────────────────────────────────────────────────────┘
```

The **ErrorDetailModal** (standalone viewer used on E2E Tests page) uses the same Split Button pattern via an `errorLogAdapter.ts` bridge.

### Full-Screen Layout

```tsx
<DialogContent className={cn(
  "flex flex-col p-0 gap-0 overflow-hidden",
  "w-full h-full max-w-full max-h-full rounded-none",          // Mobile: full screen
  "sm:max-w-[95vw] sm:w-[95vw] sm:max-h-[95vh] sm:h-[95vh] sm:rounded-lg",  // Desktop
  "lg:max-w-6xl"
)}>
```

---

*Modal structure — updated: 2026-03-31*
