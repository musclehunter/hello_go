import { useEffect, useState } from "react";
import "./App.css";

function App() {
  const [world, setWorld] = useState(null);
  const [error, setError] = useState(null);
  const [logs, setLogs] = useState([]);
  const [selectedAreaId, setSelectedAreaId] = useState(null);
  const [selectedJobTab, setSelectedJobTab] = useState(null);

  useEffect(() => {
    const url = `${import.meta.env.VITE_API_URL}/`;
    const fetchWorld = async () => {
      try {
        const res = await fetch(url);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const json = await res.json();
        setWorld(json);
        // 初回のみ選択状態を設定（既存選択がある場合は上書きしない）
        if (json.Areas && json.Areas.length > 0) {
          setSelectedAreaId(prev => (prev == null ? json.Areas[0].Id : prev));
        }
        setError(null);
      } catch (e) {
        setError(e.message);
      }
    };
    fetchWorld();
    const t = setInterval(fetchWorld, 3000); // turn_seconds=3 に合わせる
    return () => clearInterval(t);
  }, []);

  useEffect(() => {
    const url = `${import.meta.env.VITE_API_URL}/logs?n=200`;
    const fetchLogs = async () => {
      try {
        const res = await fetch(url);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const json = await res.json();
        setLogs(json);
      } catch (e) {
        // ログ取得のエラーは致命的ではないため console に出すのみ
        console.warn("fetch logs error", e);
      }
    };
    fetchLogs();
    const t = setInterval(fetchLogs, 3000);
    return () => clearInterval(t);
  }, []);

  // タブ初期化・選択維持（フック順を保つためここで宣言）
  useEffect(() => {
    if (!world) {
      setSelectedJobTab(null);
      return;
    }
    const area = world.Areas.find(a => a.Id === selectedAreaId) || world.Areas[0];
    if (!area) {
      setSelectedJobTab(null);
      return;
    }
    const personsMap = new Map(world.Persons?.map(p => [p.Id, p]) || []);
    const jobsMap = new Map(world.Jobs?.map(j => [j.Id, j]) || []);
    const names = new Set();
    (area.ResidentIDs || []).forEach(id => {
      const p = personsMap.get(id);
      if (!p) return;
      const j = jobsMap.get(p.MainJobID);
      names.add(j?.Name || "-");
    });
    const jobKeys = Array.from(names).sort((a,b)=>a.localeCompare(b));
    if (jobKeys.length === 0) {
      setSelectedJobTab(null);
      return;
    }
    setSelectedJobTab(prev => (prev && jobKeys.includes(prev) ? prev : jobKeys[0]));
  }, [world, selectedAreaId]);

  if (error) return <div style={{ padding: 16 }}>Error: {error}</div>;
  if (!world) return <div style={{ padding: 16 }}>Loading...</div>;

  // 索引を用意
  const personsById = new Map(world.Persons?.map(p => [p.Id, p]) || []);
  const jobsById = new Map(world.Jobs?.map(j => [j.Id, j]) || []);
  const selectedArea = world.Areas.find(a => a.Id === selectedAreaId) || world.Areas[0];
  const residents = (selectedArea?.ResidentIDs || [])
    .map(id => personsById.get(id))
    .filter(Boolean)
    .map(p => ({
      id: p.Id,
      name: p.Name,
      level: p.Level,
      job: jobsById.get(p.MainJobID)?.Name || "-",
      jobId: p.MainJobID,
      jobIncome: jobsById.get(p.MainJobID)?.Income ?? 0,
    }));

  // テーブル列: ジョブ名ごとに列を用意（名前順）
  const jobsSorted = [...(world.Jobs || [])].sort((a,b) => a.Name.localeCompare(b.Name));
  const jobIds = jobsSorted.map(j => j.Id);
  const jobNames = jobsSorted.map(j => j.Name);

  // エリア内のジョブ人数をカウント（jobId -> count）
  const countJobsInArea = (area) => {
    const m = new Map();
    (area.ResidentIDs || []).forEach(id => {
      const p = personsById.get(id);
      if (!p) return;
      m.set(p.MainJobID, (m.get(p.MainJobID) || 0) + 1);
    });
    return m;
  };

  // 詳細表示用: ジョブごとにグルーピング
  const residentsByJob = residents.reduce((acc, r) => {
    const key = r.job || "-";
    if (!acc[key]) acc[key] = [];
    acc[key].push(r);
    return acc;
  }, {});


  return (
    <div style={{ padding: 16 }}>
      <h1>{world.Name}</h1>
      <div className="grid">
        <div className="panel">
          <h2>Areas</h2>
          <table>
            <thead>
              <tr>
                <th>Area</th>
                <th>Population</th>
                <th>Income</th>
                {jobNames.map(n => (
                  <th key={n}>{n}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {world.Areas.map(a => {
                const counts = countJobsInArea(a);
                return (
                  <tr
                    key={a.Id}
                    className={a.Id === selectedAreaId ? "row--selected" : ""}
                    onClick={() => setSelectedAreaId(a.Id)}
                    style={{ cursor: 'pointer' }}
                  >
                    <td>{a.Name}</td>
                    <td>{a.Population}</td>
                    <td>{a.Income}</td>
                    {jobIds.map(id => (
                      <td key={id}>{counts.get(id) || 0}</td>
                    ))}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        <div className="panel">
          <h2>Area Detail</h2>
          {selectedArea ? (
            <>
              <div style={{ marginBottom: 8 }}>
                <strong>{selectedArea.Name}</strong>
                <span style={{ marginLeft: 8, color: '#6b7280' }}>ID: {selectedArea.Id}</span>
              </div>
              {residents.length === 0 ? (
                <div style={{ color: '#6b7280' }}>No residents</div>
              ) : (
                <>
                  <div className="tabs">
                    {Object.keys(residentsByJob).sort((a,b)=>a.localeCompare(b)).map(jobName => (
                      <div
                        key={jobName}
                        className={`tab ${selectedJobTab === jobName ? 'tab--active' : ''}`}
                        onClick={() => setSelectedJobTab(jobName)}
                      >
                        {jobName} <span style={{ color: '#6b7280', fontWeight: 400 }}>({residentsByJob[jobName].length})</span>
                      </div>
                    ))}
                  </div>
                  {selectedJobTab && (
                    <table>
                      <thead>
                        <tr>
                          <th>Resident</th>
                          <th>Level</th>
                          <th>Income</th>
                        </tr>
                      </thead>
                      <tbody>
                        {residentsByJob[selectedJobTab].map(r => (
                          <tr key={r.id}>
                            <td>{r.name}</td>
                            <td>{r.level}</td>
                            <td>{r.level * r.jobIncome}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                </>
              )}
            </>
          ) : (
            <div style={{ color: '#6b7280' }}>No area selected</div>
          )}
        </div>
        <div className="panel">
          <h2>Logs</h2>
          <div className="logs">
            {logs.length === 0 ? (
              <div>No logs</div>
            ) : (
              logs.map((l, idx) => (
                <div className="log-row" key={idx}>
                  <span className="log-time">{new Date(l.time).toLocaleTimeString('ja-JP', { hour12: false })}</span>
                  <span className={`badge badge--${l.type}`}>[{l.type}]</span>
                  <span>{l.message}</span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;