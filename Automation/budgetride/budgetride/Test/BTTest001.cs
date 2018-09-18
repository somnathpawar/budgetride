using System;
using budgetride.Pages;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace budgetride.Test
{
	[TestClass]
	public class BTTest001 : TestBase
	{
		[TestMethod]
		[TestProperty("Author", "Uttmesh Shukla")]
		[Description("Home Page open sucessfully")]
		public void Open_HomePage()
		{
			BaseApplicationPage baseApplicationPage = new BaseApplicationPage(Driver);
			HomePage homePage = new HomePage(Driver);

			baseApplicationPage.GoTO();
  
			homePage.LogoDisplayed();
		}
	}
}
