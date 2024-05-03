using System;
using System.Drawing;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace aNIS
{
    public partial class b2a : Form
    {
        private string link;
        private home homePage;
        public b2a(home page)
        {
            homePage = page;
            InitializeComponent();
        }

        private void label1_Click(object sender, EventArgs e)
        {
        }

        private void button1_Click(object sender, EventArgs e)
        {
            //берем только часть после b2a.kz/
            string longUrl = textBox1.Text;
            if (longUrl.StartsWith("https://"))
            {
                longUrl.Substring(8);
            }

            if (longUrl.StartsWith("b2a.kz/"))
            {
                longUrl.Substring(7);
            }
            fetch(longUrl);
        }

        private async void fetch(string longUrl)
        {
            try
            {
                //отправляем запрос на свой сервер для проверки на вирус
                using (var client = new HttpClient())
                {
                    string url = $"http://localhost:8080/virus?url={longUrl}";
                    var response = await client.GetAsync(url);

                    if (response.IsSuccessStatusCode)
                    {
                        string responseBody = await response.Content.ReadAsStringAsync();
                        Console.WriteLine(responseBody);
                        //Если он вернет ссылку, то делаем линклэйбл, иначе просто показываем результат антивирусов
                        if (responseBody.StartsWith("https://"))
                        {
                            link = responseBody;
                            label3.Text = "Перейти на virustotal";
                            linkLabel1.Text = "Открыть в Edge";
                            linkLabel1.Show();
                        }
                        else
                        {
                            label3.Text = "Опасно/всего\n"+responseBody;
                            if (responseBody[0] != '0')
                            {
                                label3.BackColor = Color.Red;
                            }
                            else
                            {
                                label3.BackColor = Color.LawnGreen;
                            }
                        }
                    }
                    else
                    {
                        MessageBox.Show("Что-то пошло не так");
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void linkLabel1_LinkClicked(object sender, LinkLabelLinkClickedEventArgs e)
        {
            System.Diagnostics.Process.Start("microsoft-edge:"+link);
        }

        private void label3_Click(object sender, EventArgs e)
        {
            
        }

        private void pictureBox2_Click(object sender, EventArgs e)
        {
            this.Hide();
            homePage.Show();
        }

        private void panel3_Paint(object sender, PaintEventArgs e)
        {
            
        }
    }
}