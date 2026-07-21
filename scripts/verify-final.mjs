import { chromium } from 'playwright';

const URL = 'https://travel.sekaishopec.com';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const ctx = await browser.newContext({
    viewport: { width: 390, height: 844 },
    deviceScaleFactor: 2,
    permissions: ['geolocation'],
    geolocation: { latitude: -2.1894, longitude: -79.8893 },
  });
  const page = await ctx.newPage();

  console.log('[1] pickup_select');
  await page.goto(URL, { waitUntil: 'networkidle', timeout: 15000 });
  await page.waitForTimeout(3500);
  
  let s = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    const sr = sheet?.getBoundingClientRect();
    const cd = sheet?.querySelector('div:nth-child(2)')?.getBoundingClientRect();
    return { sheetH: Math.round(sr?.height||0), contentH: Math.round(cd?.height||0), ws: Math.round((sr?.bottom||0)-(cd?.bottom||0)) };
  });
  console.log(`   Sheet: ${s.sheetH}px, Content: ${s.contentH}px, Whitespace: ${s.ws}px ${s.ws <= 5 ? '✓' : '✗'}`);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/final-01-pickup.png' });

  // → input
  console.log('[2] input (confirm pickup)');
  await page.getByRole('button', { name: /Confirmar ubicación/i }).click();
  await page.waitForTimeout(800);
  s = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    const sr = sheet?.getBoundingClientRect();
    const cd = sheet?.querySelector('div:nth-child(2)')?.getBoundingClientRect();
    return { sheetH: Math.round(sr?.height||0), contentH: Math.round(cd?.height||0), ws: Math.round((sr?.bottom||0)-(cd?.bottom||0)) };
  });
  console.log(`   Sheet: ${s.sheetH}px, Content: ${s.contentH}px, Whitespace: ${s.ws}px ${s.ws <= 5 ? '✓' : '✗'}`);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/final-02-input.png' });

  // Search + select destination
  console.log('[3] search destination');
  const input = page.getByRole('textbox', { name: /Buscar destino/i });
  await input.fill('Av. 9 de Octubre, Guayaquil');
  await page.waitForTimeout(1500);
  const sugCount = await page.locator('button[aria-label^="Seleccionar"]').count();
  console.log(`   Suggestions: ${sugCount}`);
  if (sugCount > 0) {
    await page.locator('button[aria-label^="Seleccionar"]').first().click();
    await page.waitForTimeout(500);
  }
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/final-03-dest-selected.png' });

  // → confirm
  console.log('[4] confirm');
  await page.getByRole('button', { name: /Buscar viaje/i }).click();
  await page.waitForTimeout(5000);
  s = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    const sr = sheet?.getBoundingClientRect();
    const cd = sheet?.querySelector('div:nth-child(2)')?.getBoundingClientRect();
    return { sheetH: Math.round(sr?.height||0), contentH: Math.round(cd?.height||0), ws: Math.round((sr?.bottom||0)-(cd?.bottom||0)) };
  });
  console.log(`   Sheet: ${s.sheetH}px, Content: ${s.contentH}px, Whitespace: ${s.ws}px ${s.ws <= 5 ? '✓' : '✗'}`);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/final-04-confirm.png' });

  // Scroll confirm to bottom
  console.log('[5] confirm scrolled');
  const sheetEl = await page.$('[aria-label="Panel de información del viaje"] > div:nth-child(2)');
  if (sheetEl) {
    await sheetEl.evaluate(el => el.scrollTo(0, el.scrollHeight));
    await page.waitForTimeout(500);
  }
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/final-05-confirm-scrolled.png' });

  // Verify pin visual centering
  const pinVis = await page.evaluate(() => {
    const wrapper = document.querySelectorAll('[style*="translate(-50%"]')[0];
    if (!wrapper) return null;
    const r = wrapper.getBoundingClientRect();
    return { x: Math.round(r.x + r.width/2), y: Math.round(r.y), centered: Math.abs((r.x + r.width/2) - window.innerWidth/2) < 10 };
  });
  console.log(`[6] Pin visual: ${JSON.stringify(pinVis)}`);

  console.log('\n[DONE]');
  await browser.close();
})();
