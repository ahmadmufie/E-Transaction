import { NextResponse } from 'next/server';
import puppeteer from 'puppeteer';

export async function POST(request: Request) {
  try {
    const data = await request.json();

    // 1. Jalankan browser tanpa layar (Headless)
    const browser = await puppeteer.launch({ headless: true });
    const page = await browser.newPage();

    // 2. Desain Template Kontrak (HTML & CSS Professional)
    const htmlContent = `
      <html>
        <head>
          <style>
            body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; padding: 50px; color: #1e293b; }
            .header { text-align: center; border-bottom: 3px solid #2563eb; padding-bottom: 20px; margin-bottom: 40px; }
            .header h1 { margin: 0; color: #0f172a; font-size: 28px; text-transform: uppercase; letter-spacing: 2px; }
            .header p { margin: 5px 0 0; color: #64748b; }
            .meta-data { display: flex; justify-content: space-between; margin-bottom: 40px; background: #f8fafc; padding: 20px; border-radius: 8px; }
            .meta-data p { margin: 5px 0; font-size: 14px; }
            table { width: 100%; border-collapse: collapse; margin-bottom: 40px; }
            th, td { border: 1px solid #cbd5e1; padding: 15px; text-align: left; }
            th { background-color: #2563eb; color: white; text-transform: uppercase; font-size: 12px; }
            td { font-size: 14px; }
            .total-row { font-weight: bold; background-color: #f1f5f9; }
            .footer { margin-top: 60px; display: flex; justify-content: space-between; }
            .sign-box { text-align: center; width: 200px; }
            .sign-space { height: 100px; border-bottom: 1px solid #0f172a; margin-bottom: 10px; }
            .watermark { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%) rotate(-45deg); font-size: 100px; color: rgba(0,0,0,0.03); z-index: -1; }
          </style>
        </head>
        <body>
          <div class="watermark">E-TRANSACT SAH</div>
          
          <div class="header">
            <h1>KONTRAK SEWA RESMI</h1>
            <p>PT KONTRAKTOR KONSTRUKSI INDONESIA</p>
          </div>

          <div class="meta-data">
            <div>
              <p><strong>ID Transaksi:</strong> ${data.transaction_id}</p>
              <p><strong>Metode Pembayaran:</strong> ${data.pay_method}</p>
            </div>
            <div style="text-align: right;">
              <p><strong>Tanggal Cetak:</strong> ${new Date().toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })}</p>
              <p><strong>Status:</strong> LUNAS / TERVERIFIKASI</p>
            </div>
          </div>

          <table>
            <thead>
              <tr>
                <th>Deskripsi Layanan / Armada</th>
                <th style="text-align: right;">Total Nilai (Rp)</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Penyewaan Alat Berat & Layanan Berdasarkan Sistem POS</td>
                <td style="text-align: right;">Rp ${Number(data.total_amount).toLocaleString('id-ID')}</td>
              </tr>
              <tr class="total-row">
                <td style="text-align: right;">GRAND TOTAL</td>
                <td style="text-align: right; color: #2563eb;">Rp ${Number(data.total_amount).toLocaleString('id-ID')}</td>
              </tr>
            </tbody>
          </table>

          <div class="footer">
            <div class="sign-box">
              <p>Pihak Penyewa,</p>
              <div class="sign-space"></div>
              <p><b>Klien / Perwakilan</b></p>
            </div>
            <div class="sign-box">
              <p>Disahkan Oleh,</p>
              <div class="sign-space"></div>
              <p><b>Manajer Operasional</b></p>
            </div>
          </div>
        </body>
      </html>
    `;

    // 3. Konversi HTML ke PDF
    await page.setContent(htmlContent, { waitUntil: 'networkidle0' });
    const pdfBuffer = await page.pdf({ 
      format: 'A4', 
      printBackground: true,
      margin: { top: '20px', bottom: '20px' }
    });
    
    await browser.close();

    // 4. Kirim PDF kembali ke peramban (browser) Kasir
    return new NextResponse(pdfBuffer, {
      status: 200,
      headers: {
        'Content-Type': 'application/pdf',
        'Content-Disposition': `attachment; filename="Kontrak-${data.transaction_id}.pdf"`,
      },
    });

  } catch (error) {
    console.error(error);
    return NextResponse.json({ error: 'Gagal membuat PDF Kontrak' }, { status: 500 });
  }
}