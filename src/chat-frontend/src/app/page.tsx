import Link from 'next/link'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const features = [
  {
    icon: '⚡',
    title: 'Real-time',
    desc: 'Messages delivered instantly via WebSocket — no polling, no delays.',
  },
  {
    icon: '🔒',
    title: 'Secure',
    desc: 'Google OAuth login. Your credentials never touch our servers.',
  },
  {
    icon: '📜',
    title: 'History',
    desc: 'Full chat history stored and searchable. Pick up right where you left off.',
  },
  {
    icon: '👥',
    title: 'Contacts',
    desc: 'Smart contact list sorted by recent activity. Find who matters, fast.',
  },
]

export default function LandingPage() {
  return (
    <main className="relative min-h-screen flex flex-col overflow-hidden">
      {/* Ambient background glow */}
      <div className="pointer-events-none fixed inset-0 z-0">
        <div
          className="absolute top-[-20%] left-[10%] w-[600px] h-[600px] rounded-full"
          style={{
            background: 'radial-gradient(circle, hsl(158 64% 52% / 0.08) 0%, transparent 70%)',
          }}
        />
        <div
          className="absolute bottom-[-10%] right-[5%] w-[500px] h-[500px] rounded-full"
          style={{
            background: 'radial-gradient(circle, hsl(220 60% 60% / 0.06) 0%, transparent 70%)',
          }}
        />
      </div>

      {/* Nav */}
      <nav className="relative z-10 flex items-center justify-between px-8 py-6 max-w-6xl mx-auto w-full">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-primary-foreground font-display font-bold text-sm">K</span>
          </div>
          <span className="font-display font-semibold text-foreground tracking-tight">Kayee</span>
        </div>
        <a
          href={`${API_URL}/google-sso`}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          Sign in →
        </a>
      </nav>

      {/* Hero */}
      <section className="relative z-10 flex-1 flex flex-col items-center justify-center text-center px-6 pt-16 pb-24">
        {/* Badge */}
        <div
          className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full border border-primary/20 bg-primary/5 text-primary text-xs font-mono mb-8"
          style={{ animation: 'fadeUp 0.5s ease forwards', opacity: 0 }}
        >
          <span className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse-dot" />
          WebSocket · Redis · Go
        </div>

        {/* Headline */}
        <h1
          className="font-display font-extrabold text-5xl md:text-7xl lg:text-8xl tracking-tight leading-[0.9] mb-6 max-w-4xl"
          style={{ animation: 'fadeUp 0.6s ease 0.1s forwards', opacity: 0 }}
        >
          Chat that{' '}
          <span
            style={{
              background: 'linear-gradient(135deg, hsl(158 64% 52%), hsl(158 64% 72%))',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
          >
            actually
          </span>
          {' '}moves.
        </h1>

        <p
          className="text-muted-foreground text-lg md:text-xl max-w-xl leading-relaxed mb-10"
          style={{ animation: 'fadeUp 0.6s ease 0.2s forwards', opacity: 0 }}
        >
          Kayee is a slang for <em>"what&apos;s up"</em> — and that&apos;s exactly what this is.
          Fast, no-frills real-time messaging built on Go and Redis.
        </p>

        {/* CTA */}
        <div
          className="flex flex-col sm:flex-row gap-3 items-center"
          style={{ animation: 'fadeUp 0.6s ease 0.3s forwards', opacity: 0 }}
        >
          <a
            href={`${API_URL}/google-sso`}
            className="group flex items-center gap-3 px-6 py-3.5 rounded-xl bg-primary text-primary-foreground font-medium text-sm hover:bg-primary/90 transition-all hover:shadow-lg hover:shadow-primary/20 hover:-translate-y-0.5"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Continue with Google
          </a>
          <a
            href="#features"
            className="px-6 py-3.5 rounded-xl border border-border text-sm text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-all"
          >
            See how it works
          </a>
        </div>

        {/* Mockup preview */}
        <div
          className="mt-16 relative w-full max-w-3xl"
          style={{ animation: 'fadeUp 0.8s ease 0.4s forwards', opacity: 0 }}
        >
          <div className="rounded-2xl border border-border bg-card overflow-hidden shadow-2xl">
            {/* Fake window bar */}
            <div className="flex items-center gap-1.5 px-4 py-3 border-b border-border bg-secondary/50">
              <div className="w-3 h-3 rounded-full bg-red-500/50" />
              <div className="w-3 h-3 rounded-full bg-yellow-500/50" />
              <div className="w-3 h-3 rounded-full bg-green-500/50" />
              <span className="ml-3 text-xs text-muted-foreground font-mono">kayee — chat</span>
            </div>
            {/* Fake chat UI */}
            <div className="flex h-48">
              {/* Sidebar */}
              <div className="w-48 border-r border-border bg-secondary/20 p-3 space-y-2">
                {['user_alpha', 'dev_team', 'kartik'].map((name, i) => (
                  <div
                    key={name}
                    className={`flex items-center gap-2 px-2 py-1.5 rounded-lg ${i === 0 ? 'bg-primary/10' : ''}`}
                  >
                    <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${i === 0 ? 'bg-primary text-primary-foreground' : 'bg-border text-muted-foreground'}`}>
                      {name[0].toUpperCase()}
                    </div>
                    <span className="text-xs text-muted-foreground truncate">{name}</span>
                  </div>
                ))}
              </div>
              {/* Messages */}
              <div className="flex-1 p-4 space-y-3 flex flex-col justify-end">
                <div className="flex justify-start">
                  <div className="bg-secondary text-secondary-foreground text-xs px-3 py-1.5 rounded-2xl rounded-tl-sm max-w-[60%]">
                    yo what&apos;s up
                  </div>
                </div>
                <div className="flex justify-end">
                  <div className="bg-primary text-primary-foreground text-xs px-3 py-1.5 rounded-2xl rounded-tr-sm max-w-[60%]">
                    nothing much, you?
                  </div>
                </div>
                <div className="flex justify-start">
                  <div className="bg-secondary text-secondary-foreground text-xs px-3 py-1.5 rounded-2xl rounded-tl-sm max-w-[60%]">
                    just shipped that feature 🚀
                  </div>
                </div>
              </div>
            </div>
          </div>
          {/* Glow under mockup */}
          <div
            className="absolute -bottom-8 left-1/2 -translate-x-1/2 w-2/3 h-16 blur-2xl"
            style={{ background: 'hsl(158 64% 52% / 0.12)' }}
          />
        </div>
      </section>

      {/* Features */}
      <section id="features" className="relative z-10 max-w-6xl mx-auto w-full px-6 pb-24">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {features.map((f, i) => (
            <div
              key={f.title}
              className="p-6 rounded-2xl border border-border bg-card hover:border-primary/20 transition-colors group"
              style={{
                animation: `fadeUp 0.6s ease ${0.1 * i + 0.5}s forwards`,
                opacity: 0,
              }}
            >
              <div className="text-2xl mb-3">{f.icon}</div>
              <h3 className="font-display font-semibold text-sm mb-1.5 text-foreground">{f.title}</h3>
              <p className="text-xs text-muted-foreground leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 border-t border-border px-8 py-6 max-w-6xl mx-auto w-full flex items-center justify-between">
        <span className="text-xs text-muted-foreground font-mono">kayee © 2026</span>
        <span className="text-xs text-muted-foreground">built with Go + Redis + Next.js</span>
      </footer>
    </main>
  )
}